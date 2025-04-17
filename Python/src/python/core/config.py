import configparser
import os
import pathlib


def read_config() -> configparser.ConfigParser:
    config = configparser.ConfigParser()
    config["Database"] = {
        "host": "localhost",
        "port": "5432",
        "username": "postgres",
        "password": "postgres",
        "database": "rosetta",
    }

    config_path = pathlib.Path(os.path.expanduser("~")).joinpath(
        ".config/rosetta/config.ini"
    )

    try:
        config.read(config_path)
    except Exception as e:
        raise FileNotFoundError(f"valid config file not found: {e}")

    return config
