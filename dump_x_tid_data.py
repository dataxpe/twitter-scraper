import json
from urllib.parse import urlparse

# pip install XClientTransaction requests beautifulsoup4

import bs4
import requests
from x_client_transaction.utils import (
    generate_headers,
    get_ondemand_file_url,
    handle_x_migration,
)
from x_client_transaction import ClientTransaction


def get_home_and_ondemand(session: requests.Session):
    """
    Get the HTML for https://x.com and the ondemand.s.*.js contents.
    """

    # Option 1: use handle_x_migration (needed if you ever hit twitter.com directly)
    # This returns a BeautifulSoup object.
    # home_page_response = handle_x_migration(session=session)

    # Option 2 (recommended for x.com): go straight to https://x.com
    home_page = session.get("https://x.com", timeout=20)
    home_page.raise_for_status()
    home_page_response = bs4.BeautifulSoup(home_page.content, "html.parser")

    # Find ondemand.s.<hex>a.js URL from the home page
    ondemand_file_url = get_ondemand_file_url(response=home_page_response)
    if not ondemand_file_url:
        raise RuntimeError("Could not find ondemand.s.*.js URL in the home page")

    ondemand_file = session.get(ondemand_file_url, timeout=20)
    ondemand_file.raise_for_status()

    # The README suggests using .text to avoid key index errors
    # (passing BeautifulSoup also works sometimes).
    ondemand_file_response = ondemand_file.text

    return home_page_response, ondemand_file_response


def extract_simple_fields(ct: ClientTransaction):
    """
    Extract only JSON-serialisable fields from ClientTransaction.__dict__
    so you can see what’s there (key, frames, etc.) without guessing names.
    """
    result = {}
    for k, v in ct.__dict__.items():
        try:
            json.dumps(v)
        except TypeError:
            # Not JSON-serialisable (e.g. BeautifulSoup, session, etc.)
            continue
        result[k] = v
    return result


def main():
    # 1) Prepare a session with headers the lib expects
    session = requests.Session()
    session.headers = generate_headers()

    # 2) Get home HTML + ondemand JS
    home_page_response, ondemand_file_response = get_home_and_ondemand(session)

    # 3) Build ClientTransaction object
    ct = ClientTransaction(
        home_page_response=home_page_response,
        ondemand_file_response=ondemand_file_response,
    )

    # 4) (Optional) show that generating a TID works
    url = "https://x.com/i/api/graphql/1VOOyvKkiI3FMmkeDNxM9A/UserByScreenName"
    method = "GET"
    path = urlparse(url).path

    tid = ct.generate_transaction_id(method=method, path=path)
    print("Example X-Client-Transaction-Id:", tid)

    # 5) Dump JSON-serialisable internals – this is what you want
    data = extract_simple_fields(ct)

    # Pretty-print to stdout
    print("\n=== JSON serialisable fields of ClientTransaction ===")
    print(json.dumps(data, indent=2, ensure_ascii=False))

    # Optionally also save to a file for your Go program
    with open("x_tid_state.json", "w", encoding="utf-8") as f:
        json.dump(data, f, indent=2, ensure_ascii=False)
    print('\nSaved JSON to x_tid_state.json')


if __name__ == "__main__":
    main()
