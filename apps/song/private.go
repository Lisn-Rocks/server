package song

// This file contains package private functions.

import (
    "net/http"
    "log"
    "database/sql"
    "path"
    "strings"
    "strconv"
    "fmt"
    "os"
    "io"
)



func parseIDFromURL(r *http.Request) (id int, err error) {
    split := strings.Split( r.URL.String(), "/" )[1:]

    // Atoi used to prevent injections. Only numbers pass!
    id, err = strconv.Atoi(split[1])

    return
}



func songExists(songid int, db *sql.DB) bool {
    var name string
    row := db.QueryRow(`SELECT name FROM songs WHERE songid=$1;`, songid)
    err := row.Scan(&name)
    return err == nil
}



func getSongExtension(songid int, db *sql.DB) (extension string, err error) {
    row := db.QueryRow(`SELECT extension FROM songs WHERE songid=$1;`, songid)
    err = row.Scan(&extension)
    return
}



func getAlbumID(songid int, db *sql.DB) (albumid int, err error) {
    row := db.QueryRow(`SELECT albumid FROM songs WHERE songid=$1;`, songid)
    err = row.Scan(&albumid)
    return
}



func getAlbumName(songid int, db *sql.DB) (name string, err error) {
    albumid, err := getAlbumID(songid, db)

    if err != nil {
        return
    }

    row := db.QueryRow(`SELECT name FROM albums WHERE albumid=$1;`, albumid)
    row.Scan(&name)
    return
}



func getAlbumExtension(albumid int, db *sql.DB) (extension string, err error) {
    row := db.QueryRow(
        `SELECT extension FROM albums WHERE albumid=$1;`, albumid,
    )
    err = row.Scan(&extension);
    return
}



func getArtistName(artistid int, db *sql.DB) (name string, err error) {
    row := db.QueryRow(`SELECT name FROM artists WHERE artistid=$1;`, artistid)
    err = row.Scan(&name)
    return
}



func serveFileFromFolder(
    w http.ResponseWriter, r *http.Request, logr *log.Logger,
    folder string, songid int, extension string,
) {
    filepath := path.Join( folder, fmt.Sprintf("%d%s", songid, extension) )

    if _, err := os.Stat(filepath); os.IsNotExist(err) {
        logr.Printf("Cannot serve. File '%s' not found in filesystem", filepath)

        w.WriteHeader(http.StatusNotFound)
        io.WriteString(w, "404 file not found")

        return
    }

    logr.Printf("Serving file %s", filepath)
    http.ServeFile(w, r, filepath)
}
