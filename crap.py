#!/usr/bin/env python3
from functools import wraps
from time import sleep
from flask import Flask, request
import random
import json

random.seed()
reliability = 3
app = Flask(__name__)
app.logger.setLevel("DEBUG")

steps = {
    "red": ["one", "two", "three", "four", "five"],
    "blue": ["six", "five", "four", "three", "two", "one"]
}

def maybe(fn):
    @wraps(fn)
    def wrapper(*args, **kwargs):
        r = random.randint(0, reliability)
        if r == 0:
            app.logger.debug("this is an error")
            raise Exception("well darn")
        return fn(*args, **kwargs)
    return wrapper

@app.route('/')
def new_job():
    app.logger.debug("new job")
    if random.randint(0, 1) == 0:
        return {"color": "blue"}
    return {"color": "red"}

@app.route('/color/<color>')
def start(color):
    app.logger.debug(color)
    return {"steps": steps[color]}

@app.route('/step/<step>')
@maybe
def number(step):
    app.logger.debug(step)
    return {"word": step}

@app.route('/done', methods=['POST'])
def done():
    app.logger.debug("done")
    content = request.json

    if content['data'] != ''.join(steps[content['color']]):
        app.logger.error(f"[FAILED] with data {content['data']}")
    else:
        app.logger.info(f"[  OK  ] {json.dumps(content)}")
    return {}

if __name__ == "__main__":
    app.run()