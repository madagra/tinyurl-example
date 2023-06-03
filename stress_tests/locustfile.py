import csv
from pathlib import Path
import random

from locust import HttpUser, task, between

BASE_URL_CLOUD = "https://mdtiny.net"

BASE_URL_LOCAL = "http://localhost:3000"

TEST_DATA_FNAME = "test_urls.csv"

get_url = lambda host, path: f"{host}/{path}"

def load_urls() -> list[str]:
    path_data = Path(".")
    res = []
    with open(path_data / TEST_DATA_FNAME, "r") as f:
        csv_reader = csv.reader(f, delimiter=",", skipinitialspace=True)
        for i, row in enumerate(csv_reader):
            
            # skip header
            if i == 0:
                continue
            
            res.append(row[0])

    return res


list_urls = load_urls()


class TinyUrlAppStress(HttpUser):

    host = BASE_URL_LOCAL
    wait_time = between(1, 10)

    @task
    def health_check(self):
        self.client.get(get_url(self.host, ""))

    @task
    def shorten_url(self):

        url = random.choice(list_urls)
        body = {    
            "url": url
        }
        self.client.post(
            get_url(self.host, "shorten"), data=body, allow_redirects=False

        )


if __name__ == "__main__":
    res = load_urls()