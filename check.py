
_map = {}
with open('output.txt','r', encoding='utf-8') as f:
    lines = f.readlines()
    for line in lines:
        tokens = line.split(" ")
        print(tokens)
        try:
            if tokens[2] == "Crawling":
                if line not in _map:
                    _map[tokens[3]] = 1
                else:
                    _map[tokens[3]] += 1
        except Exception as e:
            print(e)
    print(_map)
    