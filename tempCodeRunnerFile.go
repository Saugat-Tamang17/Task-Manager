http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
    // what user sent
    var clientUser AuthUser
    json.NewDecoder(r.Body).Decode(&clientUser)

    // fetch from database using the username user sent
    var dbUser AuthUser
    db.QueryRow("SELECT id, username, password FROM users_jwt WHERE username=$1", clientUser.Username).Scan(&dbUser.Id, &dbUser.Username, &dbUser.Password)

    // compare passwords
    err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(clientUser.Password))
    if err != nil {
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }
    token:=jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
			"userId":dbUser.Id,
			"exp":time.Now().Add(time.Hour *24).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
if err != nil {
    http.Error(w, "Error generating token", http.StatusInternalServerError)
    return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
})