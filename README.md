# stack-overflow monitor

Simple [stackoverflow](https://stackoverflow.com/) monitor (tags based)

![](https://butuzov.github.io/stackoverflow/screen.png)

## Install

```
go install github.com/butuzov/stackoverflow@latest
```


## Usage

```
stackoverflow -h=wordpress.stackexchange.com plugin-development
stackoverflow -o -h=codereview.stackexchange.com solid
stackoverflow -o go golang
```


### Options

* `-o`, `--open` - open new questions in the browser.
* `-h`, `--host` - stackoverflow host, by default `stackoverflow.com`
