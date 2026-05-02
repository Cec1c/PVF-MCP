#!/usr/bin/env python3
"""E2E test suite for PVF-MCP server. Requires pvfUtility running with a PVF loaded."""

import subprocess, json, sys, os

BINARY = os.path.join(os.path.dirname(__file__), "pvf-mcp.exe")

def run_tests():
    proc = subprocess.Popen(
        [BINARY],
        stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE
    )
    rid = [0]
    def rpc(method, params=None):
        rid[0] += 1
        msg = {"jsonrpc":"2.0","id":rid[0]}
        if method:
            msg["method"] = method
            msg["params"] = params or {}
        proc.stdin.write((json.dumps(msg)+"\n").encode())
        proc.stdin.flush()
        line = proc.stdout.readline()
        return json.loads(line)

    passed = 0
    failed = 0

    def check(name, condition, detail=""):
        nonlocal passed, failed
        if condition:
            passed += 1
            print(f"  PASS: {name}")
        else:
            failed += 1
            print(f"  FAIL: {name} - {detail}")

    print("PVF-MCP E2E Test Suite\n")

    # 1. Initialize
    print("[1] Initialize MCP session")
    init = rpc("initialize", {"protocolVersion":"2024-11-05","capabilities":{},"clientInfo":{"name":"test","version":"1.0"}})
    check("Protocol version", init['result']['protocolVersion'] == "2024-11-05")
    check("Server name", init['result']['serverInfo']['name'] == "pvf-mcp")
    check("Tools capability", init['result']['capabilities']['tools']['listChanged'] == True)

    # 2. List tools
    print("\n[2] List tools")
    tl = rpc("tools/list")
    tools = [t['name'] for t in tl['result']['tools']]
    check("26 tools registered", len(tools) == 26, f"Got {len(tools)}")
    essential = ["get_version","search_pvf","get_file_data","get_file_content","import_file","save_pvf","serialize_file_data"]
    for t in essential:
        check(f"Essential tool: {t}", t in tools)

    # 3. get_version
    print("\n[3] Connection test")
    ver = rpc("tools/call", {"name":"get_version","arguments":{}})
    v = ver['result']['content'][0]['text']
    check("get_version returns string", len(v) > 0 and "." in v, v)

    # 4. get_loaded_pvf_path
    path = rpc("tools/call", {"name":"get_loaded_pvf_path","arguments":{}})
    pvp = path['result']['content'][0]['text']
    check("PVF loaded", len(pvp) > 0 and "Script.pvf" in pvp, pvp)

    # 5. get_pvf_root_directory
    print("\n[4] Directory listing")
    root = rpc("tools/call", {"name":"get_pvf_root_directory","arguments":{}})
    dirs = json.loads(root['result']['content'][0]['text'])
    check("Root dirs returned", len(dirs) > 10, f"{len(dirs)} dirs")
    check("Has equipment", "equipment" in dirs)
    check("Has stackable", "stackable" in dirs)

    # 6. search_pvf
    print("\n[5] Search")
    sr = rpc("tools/call", {"name":"search_pvf","arguments":{"keyword":"SP+20"}})
    results = json.loads(sr['result']['content'][0]['text'])
    check("Found SP+20 files", len(results) >= 1)
    check("Contains book_skill2", any("book_skill2" in r for r in results))

    # 7. get_item_info
    print("\n[6] Item info")
    info = rpc("tools/call", {"name":"get_item_info","arguments":{"file_path":"stackable/book_skill2.stk"}})
    item = json.loads(info['result']['content'][0]['text'])
    check("Item name correct", item['ItemName'] == "SP+20技能书")
    check("Item code correct", item['ItemCode'] == 1038)

    # 8. get_file_data (structured)
    print("\n[7] Structured file data")
    fd = rpc("tools/call", {"name":"get_file_data","arguments":{"file_path":"stackable/book_skill2.stk"}})
    data = json.loads(fd['result']['content'][0]['text'])
    check("Has sections", len(data) > 10)
    price_node = None
    for node in data:
        if node.get('SectionName') == '[price]':
            price_node = node
            break
    check("Found [price] section", price_node is not None)
    if price_node:
        check("Price has value", len(price_node['Children']) >= 1)

    # 9. get_file_content (raw)
    print("\n[8] Raw file content")
    raw = rpc("tools/call", {"name":"get_file_content","arguments":{"file_path":"stackable/book_skill2.stk"}})
    content = raw['result']['content'][0]['text']
    check("Raw starts with #PVF_File", content.startswith("#PVF_File"))
    check("Contains price data", "133322" in content or "132234" in content or "132999" in content)

    # 10. file_exists
    print("\n[9] File existence")
    ex1 = rpc("tools/call", {"name":"file_exists","arguments":{"file_path":"stackable/book_skill2.stk"}})
    check("Existing file found", ex1['result']['content'][0]['text'] == "true")
    ex2 = rpc("tools/call", {"name":"file_exists","arguments":{"file_path":"nonexistent/xyz.file"}})
    check("Missing file not found", ex2['result']['content'][0]['text'] == "false")

    # 11. item_code_to_file_info
    print("\n[10] Code lookup")
    code = rpc("tools/call", {"name":"item_code_to_file_info","arguments":{"lst_names":"stackable","item_code":1038}})
    cinfo = json.loads(code['result']['content'][0]['text'])
    check("Code→file path", "book_skill2" in cinfo.get('FilePath',''))
    check("Code→item name", "SP+20" in cinfo.get('ItemName',''))

    # 12. get_all_lst_file_list
    print("\n[11] LST files")
    lst = rpc("tools/call", {"name":"get_all_lst_file_list","arguments":{}})
    lsts = json.loads(lst['result']['content'][0]['text'])
    check("Has equipment.lst", any("equipment.lst" in l for l in lsts))
    check("Has stackable.lst", any("stackable.lst" in l for l in lsts))

    # 13. get_lst_file_info
    print("\n[12] LST info")
    lsti = rpc("tools/call", {"name":"get_lst_file_info","arguments":{"file_path":"stackable/stackable.lst"}})
    lst_data = json.loads(lsti['result']['content'][0]['text'])
    check("Stackable lst has entries", len(lst_data) > 0)
    check("Contains item 1038", "1038" in lst_data)

    # 14. Serialize roundtrip
    print("\n[13] Serialize roundtrip")
    ser = rpc("tools/call", {"name":"serialize_file_data","arguments":{"nodes":json.dumps(data)}})
    serialized = ser['result']['content'][0]['text']
    check("Serialized has #PVF_File", serialized.startswith("#PVF_File"))
    check("Serialized has [price]", "[price]" in serialized)

    # 15. get_selected_files (GUI)
    print("\n[14] GUI interaction")
    sf = rpc("tools/call", {"name":"get_selected_files","arguments":{}})
    sel = json.loads(sf['result']['content'][0]['text'])
    check("Get selected files works", isinstance(sel, list))

    # 16. get_string_table
    print("\n[15] String table")
    st = rpc("tools/call", {"name":"get_string_table","arguments":{}})
    st_data = json.loads(st['result']['content'][0]['text'])
    check("String table has entries", len(st_data) > 100, f"{len(st_data)} entries")

    # Summary
    total = passed + failed
    print(f"\n{'='*50}")
    print(f"Results: {passed}/{total} passed, {failed} failed")
    print(f"{'='*50}")

    proc.terminate()
    return failed == 0

if __name__ == "__main__":
    success = run_tests()
    sys.exit(0 if success else 1)
