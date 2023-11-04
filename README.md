# TinyURL sample project 

Sample project to experiment with the system design with a TinyURL application

## Test the application locally

Use the following Python code to test the application locally:

```python
import requests

url = "replace_with_valid_url"

# shorten
myobj = {'url': url}
resp = requests.post(get_url(BASE_URL_LOCAL, "shorten"), json=myobj)

# redirect
short_url = resp.text
requests.get(x.text, allow_redirects=True)
```

## Run stress tests

Stress tests are implemented using the `locust` open-source library. In order to run the stress tests, get into the `stress_tests` directory
and execute the following command:

```bash
python -m venv .venv
source .venv/bin/activate
python -m pip install -r requirements.txt

# increase temporarily the maximum number of
# file descriptors allowed
ulimit -n 50000

# run the Locust GUI on the localhost
locust

# if testing locally start the application
cd ../tinyurl
go build -c app && ./app
```

After running these commands, the `locust` interface will be available on `http://localhost:8089`. Use the instructions on the GUI to start the stress testing and
select `http://localhost:3000` as a target if you want to test the application locally.
