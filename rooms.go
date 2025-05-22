package main

type Room struct {
	Name    string
	Body    string
	Links   string
	Image   []byte
	Style   string
	User_id int
}

func getRoom(user string) (Room, error) {
	var room Room

	err := db.QueryRow("select name, body, links, image, style, user_id from rooms where lower(name) = lower($1)", user).Scan(&room.Name, &room.Body, &room.Links, &room.Image, &room.Style, &room.User_id)

	if err != nil {
		return room, err
	}

	return room, nil
}
