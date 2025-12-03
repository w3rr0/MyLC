# Serwer dla stanowisk IAESTE
Poniższa instrukcja pomoże Ci skonfigurować środowisko,
przygotować bazę danych i uruchomić aplikację lokalnie.

## 📋 Wymagania wstępne
Zanim zaczniesz, upewnij się, że masz zainstalowane następujące narzędzia:

- `Go`: w wersji 1.25 lub nowszej.

- `PostgreSQL`: w wersji 17.

- `Git`: do pobrania repozytorium.

## 📲 Jak zainstalować

#### Linux / macOS Najłatwiej użyć menedżera pakietów

```bash
# macOS
brew install go postgresql

# Linux (Ubuntu/Debian)
sudo apt update
sudo apt install golang postgresql
```

#### Alternatywnie pobierz instalatory ze stron oficjalnych:
[Go Download](https://go.dev/doc/install) |
[PostgreSQL Download](https://www.postgresql.org/download/)

#### Windows Pobierz i uruchom instalatory:

[Go 1.25 Installer](https://go.dev/doc/install)

[PostgreSQL 17 Installer](https://www.postgresql.org/download/)

## 📥 Pobranie Repozytorium

Zforkuj to repozytorium na swoje konto GitHub,
a następnie sklonuj je lokalnie:

```bash
# Za pomocą ssh
git clone git@github.com:<YOUR_NICK>/IAESTE_stands_server.git
# Lub https
git clone https://github.com/<YOUR_NICK>/IAESTE_stands_server.git
```
Gałąź główna `main` to zawsze najnowsza stabilna wersja kodu.
Gdy wychodzi nowa wersja, zaktualizuj swoją kopię kodu za pomocą przycisku `Sync fork`
w repo na swoim koncie.

## 🗄️ Konfiguracja Bazy Danych

Projekt wymaga bazy danych PostgreSQL.
Plik ze strukturą bazy znajduje się w `schema.sql`.

Upewnij się, że serwer PostgreSQL jest uruchomiony.

Stwórz nową bazę danych o nazwie `test` (nazwa zostanie zmieniona,
gdy wymyślę jak cały projekt ma się nazywać) i zaimportuj schemat.

#### Linux / macOS (Terminal)
```bash
# 1. Stwórz bazę danych
createdb test

# 2. Zaimportuj schemat
psql -d test -f schema.sql
```

#### Windows (PowerShell lub CMD)
```powershell
# 1. Stwórz bazę danych (możesz też użyć pgAdmin)
createdb -U postgres test

# 2. Zaimportuj schemat
psql -U postgres -d moj_projekt_db -f schema.sql
```

## 🚀 Uruchomienie Projektu

Gdy masz już zainstalowane Go i przygotowaną bazę danych,
wykonaj następujące kroki:

Zmień nazwę projektu na `go_server` - ta nazwa również zostanie zmieniona po ustaleniu nazwy projektu.

1. Pobierz zależności

Otwórz terminal w folderze projektu i uruchom:

```bash
go mod download
```

2. Uruchom serwer

Punkt wejściowy aplikacji znajduje się w `cmd/server/main.go`.

#### Linux / macOS / Windows Komenda jest identyczna dla wszystkich systemów:
```bash
go run cmd/server/main.go
```

Jeśli wszystko poszło zgodnie z planem, powinieneś zobaczyć w konsoli informację,
że serwer wystartował (Server is running at :8080).

## 🛠️ Rozwiązywanie problemów

- Do pomocy z setupem projektu można pisać do mnie DM
- W przypadku znalezienia jakiegoś błędu w programie proszę o stworzenie issue
w tym repo, najlepiej z opisem błędu oraz krokami do jego odtworzenia.
W pierwszej kolejności jednak sprawdź, czy ktoś inny nie zgłosił już tego problemu
oraz czy nie został on już naprawiony w najnowszej wersji kodu.

## 📄 Uwaga

Ten projekt jest w fazie wczesnego rozwoju i może ulec znacznym zmianom.
Jeśli zajmujesz się jego rozwojem, zalecam regularnie zaglądać do tego repo.