# Welcome mailguard

Automated mail-bot written in go language.  Dynamic behaviour loading an yml file and evaluating it using a rule-engine.

The first verison is in process...


## How to build

Create and fill conf.json and rules.yml files (check examples folder).

Execute build script.

    ./scripts/build.sh


## How it works

![High level diagram](https://github.com/IberLabs/mailguard/blob/master/doc/mailguard.png?raw=true)

1. Connect to POP3/IMAP mailbox.
2. Load dynamic ruleset from configuration file.
3. Move trough the e-mail list evaluating dynamic rules.
4. Perform configured action from a set of available actions: send e-mail, trigger webhook, integrations.


## Collaboration

This project is on an early stage but any collaboration would be more than appreciated.