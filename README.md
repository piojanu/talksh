# talksh

`talksh` is a minimal, Unix-friendly command-line tool for running LLM prompts on data. It reads from `stdin`, calls an LLM, and prints the result to `stdout`.

It's designed to be a simple, powerful building block in your shell, letting you compose it with other tools like `cat`, `grep`, `jq`, `pv`, etc.

For example, to count the number of posts about AI in your blog:
```bash
cat blog-posts.jsonl | jq '.title' | \
  talksh map --prompt 'Is this post about AI? Answer with exactly YES or NO, nothing else.\n\n{{}}' | \
  grep 'YES' | wc -l
```

-----

## Installation

There are two main ways to install `talksh`.

### Option 1: GitHub Actions Artifacts (Recommended)

Every push builds fresh binaries. Grab the latest run:

1. Visit https://github.com/piojanu/talksh/actions
2. Click the most recent ✅ run on `main` (or the tag you care about)
3. Download the `talksh-<os>-<arch>` artifact you need
4. Rename and move the `talksh` binary to a directory in your `$PATH`.

**For non-root users (Recommended):**

A common user-specific location is `$HOME/.local/bin`.

```bash
# Ensure the directory exists
mkdir -p $HOME/.local/bin

# Move the binary
mv talksh-<os>-<arch> $HOME/.local/bin/talksh
```

*(Note: Make sure `$HOME/.local/bin` is in your `$PATH` by adding `export PATH="$HOME/.local/bin:$PATH"` to your `~/.bashrc` or `~/.zshrc` if it isn't already.)*

**For root users (System-wide):**

You can make `talksh` available to all users by placing it in `/usr/local/bin`.

```bash
sudo mv talksh-<os>-<arch> /usr/local/bin/talksh
```

### Option 2: With `go install` (For Go Users)

If you have the [Go toolchain](https://go.dev/doc/install) installed, you can install `talksh` with a single command:

```bash
go install github.com/piojanu/talksh@latest
```

This will compile and place the `talksh` binary in your `$GOPATH/bin` directory. Make sure this directory is in your system's `$PATH`.

-----

## Configuration

`talksh` is configured via a file or environment variables, loaded in this order of priority:

1.  **Default Config File**: `$HOME/.talksh.yaml`
2.  **Custom Config File**: (Specified with `--config /path/to/config.yaml`)
3.  **Environment Variables**

### 1\. Config File (Recommended)

By default, `talksh` looks for a config file at `$HOME/.talksh.yaml`.
Create this file with the following structure:

```yaml
# $HOME/.talksh.yaml
api:
  base_url: "https://api.openai.com/v1"
  key: "sk-your-key-here"
  model: "gpt-4o"
```

You can also specify a custom config file path using the global `--config` flag:

```bash
cat ... | talksh --config /path/to/my-config.yaml map ...
```

### 2\. Environment Variables

You can override any config file setting with environment variables. `talksh` automatically maps `.` in the keys to `_` (e.g., `api.key` becomes `API_KEY`).

  * `API_KEY` overrides `api.key`
  * `API_BASE_URL` overrides `api.base_url`
  * `API_MODEL` overrides `api.model`

**Example:**

```bash
# This key will be used for this session, ignoring the config file's key
export API_KEY="sk-other-key-for-this-session"
cat ... | talksh map ...
```

-----

## Commands

`talksh` has two simple commands that represent the "map" and "reduce" patterns.

### `talksh map`

A **1-to-1** operation. It runs your prompt **for each line** of input it receives.

  * 100 lines of input = 100 LLM calls.
  * **Use for**: Transforming, extracting, or analyzing many small items (like files, log lines, or JSON objects).

### `talksh reduce`

An **N-to-1** operation. It reads the **entire input** as a single block of text and runs your prompt **once**.

  * 100 lines of input = 1 LLM call.
  * **Use for**: Summarizing, synthesizing, or getting a single answer about a large amount of text.

-----

## How Prompts Work: The `{{}}` Placeholder

Your prompt tells `talksh` *what to do* with the content you pipe in. Inside your prompt, you **must** use the `{{}}` placeholder to mark where the input content should go.

  * For `talksh map`, `{{}}` is replaced by **each line**, one at a time.
  * For `talksh reduce`, `{{}}` is replaced by the **entire text block** at once.

You can provide a prompt in two ways:

1.  As a string: `--prompt "Summarize this: {{}}"`
2.  From a file: `--prompt-file ./my-prompts/summarize.txt`

-----

## Examples (From Simple to Powerful)

Here are common use cases, starting with the basics.

### Example 1: The Basics ("Hello World")

Pipe a line of text into `talksh map`. It will replace `{{}}` with "Hello World" and send it to the LLM.

```bash
echo "Hello World" | talksh map --prompt 'Say "{{}}"'
```

Pipe multiple lines of text into `talksh map`. It will replace `{{}}` with each line of text and send it to the LLM one by one.

```bash
printf 'foo\nbar\n' | talksh map --prompt 'Say "{{}}"'
```

### Example 2: Using `talksh` as a Shell Helper

If your prompt doesn't need any input text, you can use `echo ""` to pipe in a single empty line to trigger the command. This is great for asking one-off questions.

```bash
echo "" | talksh map --prompt "How do I use 'xargs' to cat the content of files listed in a text file?"
```

### Example 3: Summarizing Multiple Files (`reduce`)

To summarize *all* your Markdown files into one answer, `cat` them all together and pipe them into `talksh reduce`. `{{}}` will be replaced by the combined text of all files.

```bash
cat *.md | talksh reduce --prompt "Summarize all the following text into 5 bullet points: {{}}"
```

### Example 4: Processing Structured Data (JSONL)

This is a primary use case. `talksh map` is perfect for JSONL (JSON Lines) files, where each line is a separate JSON object. `talksh map` sees each line (object) as one item.

```bash
cat data.jsonl | talksh map --prompt "Extract the 'user' and 'timestamp' from this JSON: {{}}"
```

### Example 5: Processing Complex JSON (with `jq`)

What if you have one giant, pretty-printed JSON file (`data.json`) that is an array?

```json
[
  {
    "user": "alice",
    "post": "Hello world"
  },
  {
    "user": "bob",
    "post": "Hi there"
  }
]
```

Use `jq` to "explode" the array and "compact" (with `-c`) each object onto a single line. `talksh map` can then process them one by one.

```bash
# jq -c '.[]' turns the file above into a stream of single-line JSON objects:
# {"user":"alice","post":"Hello world"}
# {"user":"bob","post":"Hi there"}

cat data.json | jq -c '.[]' | talksh map --prompt "Who wrote this post? {{}}"
```

### Example 6: Running Multiple Prompts (with a Shell Loop)

`talksh` only accepts one prompt at a time. To run multiple prompts over the same data, use a shell `for` loop.

Let's say you have a directory `~/my-prompts/` with:

  * `01-summarize.txt` (contains "Summarize this: {{}}")
  * `02-extract-topics.txt` (contains "Extract 3 main topics from this: {{}}")

You can run them all against your files like this:

```bash
# This loop runs once for each prompt file
for p in ~/my-prompts/*.txt; do
  echo "--- Running prompt: $p ---" >&2
  cat *.md | talksh reduce --prompt-file "$p"
done
```

*(Note: `>&2` prints the "Running prompt" message to `stderr`, so it doesn't get mixed up with your real output.)*

### Example 7: The Map-Reduce Pipeline (Advanced)

This is the most powerful use case. You can pipe the output of `talksh map` directly into `talksh reduce` for a two-stage analysis.

Let's say you have `logs.jsonl` full of server logs.

**Step 1. `map`**: First, "map" over every log line. Ask the LLM to identify a specific error. If the error isn't found, it should output a consistent string like `NO_PROBLEM`.

**Step 2. `reduce`**: Second, "reduce" the *results* from the map step. Pipe the output of `map` into `reduce` and ask a final question about the collected results.

```bash
# Define the prompts
MAP_PROMPT="Does this log line contain a 'Database Connection Error'? If yes, describe the error. If no, just say 'NO_PROBLEM'. Here is the log: {{}}"
REDUCE_PROMPT="You will receive a list of error descriptions. Count the total number of errors (ignore lines that say 'NO_PROBLEM') and provide 3 representative examples. Here is the list: {{}}"

# Run the pipeline
cat logs.jsonl | talksh map --prompt "$MAP_PROMPT" | talksh reduce --prompt "$REDUCE_PROMPT"
```

This command first runs `talksh map` hundreds of times (once per log line), and then runs `talksh reduce` *once* on the collected output of the first stage.

-----

## Bonus: Tracking Progress

LLM calls are slow, especially for `map` commands. You can track progress by piping your data through `pv` (Pipe Viewer), a standard Unix tool (`brew install pv` or `apt-get install pv`).

`pv`'s behavior changes depending on whether it can guess the *total size* of the stream.

### Method 1: The Full Progress Bar (Known Size)

If you pipe a **single file** into `pv`, it can read the file's size and show you a full percentage-based progress bar. This is the simplest case.

```bash
# pv knows the total size of logs.jsonl and shows a real-time %
cat logs.jsonl | pv | talksh map --prompt "Is this an error log? {{}}"
```

**Terminal Output:**

```
[=======>             ] 35%  1.2MiB/s
```

### Method 2: The Live Line Counter (Unknown Size)

This is the most common case for pipelines. If `pv` receives a stream of *unknown* size (like the output from `jq` or `cat *`), it can't show a percentage.

Instead, you can use the `-l` (line count) flag. This shows a live-updating counter, which is perfect for tracking `talksh map`.

```bash
# pv -l counts the lines coming out of jq
cat data.json | jq -c '.[]' | pv -l | talksh map --prompt "Processing... {{}}"
```

This also applies to piping multiple files:

```bash
# pv -l counts the total lines from all .md files
cat *.md | pv -l | talksh reduce --prompt "Summarize this... {{}}"
```

**Terminal Output:**

```
[ 145 ]
```

*(This shows 145 lines have been processed so far)*

### Method 3: Forcing a Progress Bar (The Temp File Trick)

What if you have a complex stream (like from `jq` or `cat *`) but you *really* want a percentage bar?

You must use a two-step "temp file" trick. This forces the stream to be fully processed first, so `pv` can read it as a single file (like in Method 1).

**Step 1. Run your slow `jq` process and save the result.**

```bash
cat data.json | jq -c '.[]' > /tmp/processed.jsonl
```

**Step 2. Run `talksh` on that file with `pv`.**
Now `pv` can see the file's total size and show a percentage.

```bash
cat /tmp/processed.jsonl | pv | talksh map --prompt "Processing... {{}}"
```

**Terminal Output:**

```
[=======>             ] 35%  1.2MiB/s
```