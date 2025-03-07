echo -n "Enter migration name: "
read -r name
goose create $name sql