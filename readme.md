# Flakey and Crap.py

These are just a couple of very short applications (scripts almost?) that I used for experimenting with temporal failure and retry modes. Keeping them around just for reference.

## Crap.py

Run with python3. Don't forget to install flask. This just exposes a couple of http endpoints for tests. Apply the "@maybe" decorator to any of the functions to make it error some percentage of the time.

#### GET /

Returns a color: 

    {"color": <color>}

    <color> is either "red" or "blue".

#### GET /color/<color>

Returns a list of "steps".

    {"steps": [<step>, ...]}

    <step> is a string, e.g. "one"

#### GET /step/<step>

Returns a "word" that equals whatever was passed in.

    {"word": <step>}

#### POST /done

Accepts the following json:

    {
        "color": <color>,
        "data": <data>
    }

    <color> is either "red" or "blue".
    <data> is a string.

Returns nothing

If the data received is equal to the concatenated step strings associated with the color, it logs a successful entry. Otherwise it logs a failure.

## Flakey

Flakey is a simple workflow that just calls all of the endpoints of Crap.py. Activities in the workflow are as follows:

- First, call `/` to get a color.
- Then call `/color/<color>` to get a list of steps.
- Call `/step/<step>` once for each step and concatenate the resulting strings.
- Post the results to `/done`.