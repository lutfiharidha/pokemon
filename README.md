<h3 align="center">Pokémon Battle Royale</h3>

<div align="center">

[![Status](https://img.shields.io/badge/status-active-success.svg)]()
[![GitHub Issues](https://img.shields.io/github/issues-pr/kylelobo/The-Documentation-Compendium.svg)](https://github.com/lutfiharidha/sequis-test/issues)
[![GitHub Pull Requests](https://img.shields.io/github/issues-pr/kylelobo/The-Documentation-Compendium.svg)](https://github.com/lutfiharidha/sequis-test/pulls)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](/LICENSE)

</div>

## 📝 Table of Contents

- [About](#about)
- [Getting Started](#getting_started)
- [Usage](#usage)
- [Built Using](#built_using)
- [Authors](#authors)

## 🧐 About <a name = "about"></a>

First of all, I apologize because I don't know the flow of how to play pokemon because I don't follow the game or the pokemon movie itself.

This project is a pokemon battle royale game where 5 pokemon will compete at once, in this project also records every match log and can be retrieved at any time and within a certain period. 

## 🏁 Getting Started <a name = "getting_started"></a>

### Installing

A step by step series of examples that tell you how to get a development env running.

- Rename .env.example file to .env
- Configuration database in .env file
```
DB_HOST= //Database Host
DB_PORT= //Database Port
DB_NAME= //Database Name
DB_USERNAME= //Database Username
DB_PASSWORD= //Database Password
```
Don't worry, you don't need to create a table for this application because this application will automatically create a table by itself. Just make sure your database is connected correctly.

## 🎈 Usage <a name="usage"></a>

```
Pokémon Battle Royale

Usage:
make battle:
This command to start the Pokémon Battle Royale.

make pokedex:
This command is to view pokedex data that specifies how many pokemon are in pokedex.

make pokedex-s:
This command is to see all pokedex data .

make log-i:
This command for creating an interval log from this day
there required to input the time interval for an example
the input 1d which means 1 day for the month change "d" to "m" and the year use "y".

make log-d:
This command for creating a specific date log.

make log-si:
This command for creating a specific interval log.

make help:
This command displays a list of commands that can be executed
```

## ⛏️ Built Using <a name = "built_using"></a>

- [Golang](https://go.dev/) - Server Environment
- [NutsDB](https://github.com/nutsdb/nutsdb) - Caching Database
- [MySQL](https://www.mysql.com/) - Database

## ✍️ Authors <a name = "authors"></a>

- [@lutfiharidha](https://github.com/lutfiharidha) 

