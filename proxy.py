import requests
from concurrent.futures import ThreadPoolExecutor, as_completed

def check_proxy(proxy_line):
    proxy = proxy_line.strip()
    proxies = {
        'http': f'http://{proxy}',
        'https': f'http://{proxy}'
    }

    try:
        res = requests.get("https://ipinfo.io/json", proxies=proxies, timeout=5)
        if res.status_code == 200:
            print(f"✅ WORKING: {proxy}")
            return proxy
        else:
            print(f"❌ BAD STATUS CODE: {proxy}")
    except Exception as e:
        print(f"❌ FAILED: {proxy} | Error: {e}")
    return None

def func():
    with open('./proxies.txt', 'r', encoding='utf-8') as f:
        lines = f.readlines()

    with ThreadPoolExecutor(max_workers=10) as executor:
        futures = []
        for line in lines:
            futures.append(executor.submit(check_proxy, line))

        for future in as_completed(futures):
            _ = future.result()  # optional: do something with the working proxies

if __name__ == "__main__":
    func()
