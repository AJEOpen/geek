func FindNickName(playerid int64) (nickName string, err error) {
	err := db.QueryRow("select nick_name from users where id = ?", playerid).Scan(&nickName)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", errors.Wrapf(err, "FindNickName(playerid:%d)", playerid)
		} else {
			...
		}
	}
	
	return nickName, nil
}
