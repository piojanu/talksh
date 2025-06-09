# talksh

`talksh` translates your natural language descriptions into shell commands. It helps you discover and use shell commands without needing to memorize them.

## Usage

```bash
$ talksh ask "find all .txt files modified today"
# talksh suggests:
> find . -type f -name '*.txt' -mtime 0
```

## Installation

Make sure you have Go installed (>= 1.18).

```bash
go install github.com/piojanu/talksh@latest
```

This will install the `talksh` binary in your `$GOPATH/bin`.

## LLM Provider

`talksh` uses an LLM to generate shell commands based on your natural language description.

### Ollama Setup (Default)

By default, `talksh` is configured to use the `gemma3:12b-it-qat` model hosted by a local Ollama instance at the default address `http://localhost:11434/v1`.

If you don't have Ollama installed:

- For Linux, you can install it with:

  ```bash
  curl -fsSL https://ollama.com/install.sh | sh
  ```

- For other operating systems, please visit [https://ollama.com/download](https://ollama.com/download).

### Using Other LLM Providers

You can configure `talksh` to use other LLM providers, such as NVIDIA's [build.nvidia.com](https://build.nvidia.com).
To do this, create or edit the configuration file at `$HOME/.talksh.yaml`:

```yaml
api:
  base_url: "https://integrate.api.nvidia.com/v1"
  key: "nvapi-..." # Replace with your NVIDIA API key
  model: "nvidia/llama-3.1-nemotron-70b-instruct"
```

## Customization

Beyond the LLM provider, you can further customize `talksh`'s behavior via the `$HOME/.talksh.yaml` configuration file.

### Target Shell

- **`assistant.shell`**: Specifies the target shell for which `talksh` should generate commands.
  - Default: `"bash"`
  - To generate commands for zsh:

    ```yaml
    assistant:
      shell: "zsh"
    ```

### LLM Guidance

- **`assistant.system_msg_tmpl`**: Defines the template for the system message sent to the LLM, guiding its response format. It includes a `%s` placeholder which `talksh` replaces with the value of `assistant.shell`.
  - Default: `"When you are asked to do something, first think step by step and then answer with a %s one-liner in the code block."`
  - To modify the template:

    ```yaml
    assistant:
      system_msg_tmpl: "You are a shell command expert. Provide the best %s command for the user's request."
    ```

### API Request Timeout

- **`api.timeout`**: Sets the timeout in seconds for requests to the LLM provider.
  - Default: `"30"`
  - To set a 60-second timeout:

    ```yaml
    api:
      timeout: "60"
    ```

### Using Environment Variables

All configuration parameters mentioned above can be set using environment variables. This is useful to e.g. quickly switch between different LLMs or set the API key. The mapping is as follows:

- Keys are converted to uppercase.
- Periods (`.`) in the YAML key are replaced with underscores (`_`).

For example, to set the API base URL and the assistant shell, you would use:

```bash
export API_BASE_URL="https://your.api.server/v1"
export ASSISTANT_SHELL="zsh"
```

**Note:** Environment variables will override values set in the `$HOME/.talksh.yaml` file, which in turn override the application's default values.

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
