package main

type Room struct {
	Name    string
	Body    string
	Links   string
	Style   int
	User_id int
}

func getRoom(user string) (Room, error) {
	var room Room

	err := db.QueryRow("select name, body, links, style, user_id from rooms where name = $1", user).Scan(&room.Name, &room.Body, &room.Links, &room.Style, &room.User_id)

	if err != nil {
		return room, err
	}

	return room, nil
}
