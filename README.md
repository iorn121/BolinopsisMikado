# Rust RPN (Reverse Polish Notation) CLI

This tool allows you to calculate the formula of Reverse Polish Notation written in standard input or in a text file.

# Features

If you install Cargo, Rust's build system and package manager environment, you can perform the calculation by simply typing the `rpncalc` command.

In addition, there is an optional feature that display the calculation process.

# Requirement

[environment]

- Windows 10
- cargo = "1.60.0"

[dependencies]

- anyhow = "1.0.58"
- clap = "3.2.14"

# Installation

You have to install Cargo from [The Rust community’s crate registry](https://crates.io/) or [rustup](https://www.rust-lang.org/tools/install).

```bash
git clone https://github.com/iorn121/rust_RPN_cli

cargo add clap
cargo add anyhow

cd rust_RPN_cli
cargo install --path .
```

# Usage

**stdin**

```bash
rpncalc

3 4 +
# 7 will be output
```

**text file**

```bash
echo "4 2 -" > rpn.txt

rpncalc rpn.txt
# 2 will be output
```

**option**

```bash
rpncalc -v

2 7 + 3 - 7 * 3 /
# ["/", "3", "*", "7", "-", "3", "+", "7"] [2]
# ["/", "3", "*", "7", "-", "3", "+"] [2, 7]
# ["/", "3", "*", "7", "-", "3"] [9]
# ["/", "3", "*", "7", "-"] [9, 3]
# ["/", "3", "*", "7"] [6]
# ["/", "3", "*"] [6, 7]
# ["/", "3"] [42]
# ["/"] [42, 3]
# [] [14]
# 14
# The above will be output
```

# Note

This tool was created in Windows environment, so I have not verified that this will work on Mac or Linux.

# Author

- iorn121(Io)
- github: @iorn121
- twitter: @121Tkn
- email: wmt.tkn.121@gmail.com
- portfolio: https://iorn121.github.io/

# License

This tool is under [MIT license](https://en.wikipedia.org/wiki/MIT_License).
