import csv
from pathlib import Path
import random
import glob
import os

from locust import HttpUser, task, between

BASE_URL_PROD = "https://mdtiny.net"

BASE_URL_LOCAL = "http://localhost:3000"

get_url = lambda host, path: f"{host}/{path}"


def load_urls() -> list[str]:
    path_data = Path(".") / "data"
    res = []
    csv_files = glob.glob(os.path.join(path_data, "*.csv"))
    for file in csv_files:
        with open(file, "r") as f:
            csv_reader = csv.reader(f, delimiter=",", skipinitialspace=True)
            for i, row in enumerate(csv_reader):
                # skip header
                if i == 0:
                    continue

                res.append(row[0])

    return res


list_urls = load_urls()

test_url_db = {}


class TinyUrlAppStress(HttpUser):
    host = BASE_URL_LOCAL
    wait_time = between(1, 10)

    @task(10)
    def health_check(self) -> None:
        self.client.get(get_url(self.host, ""))

    @task(100)
    def shorten_url(self) -> None:
        url = random.choice(list_urls)
        body = {"url": url}
        resp = self.client.post(
            get_url(self.host, "shorten"), data=body, allow_redirects=False
        )
        test_url_db[url] = resp.text

    @task(100)
    def redirect(self) -> None:
        try:
            short_url = random.choice(list(test_url_db.values()))
            self.client.get(short_url, allow_redirects=False)
        except IndexError:
            pass


if __name__ == "__main__":
    res = load_urls()

    import requests

    url = random.choice(res)
    myobj = {"url": url}
    x = requests.post(get_url(BASE_URL_LOCAL, "shorten"), json=myobj)
    print(x.text)
    x = requests.get(x.text, allow_redirects=True)
    print(x.status_code)
