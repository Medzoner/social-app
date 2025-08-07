package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"social-app/internal/config"
	"social-app/internal/connector"
	"social-app/internal/models"
)

var mediaTypes = map[string][]struct {
	Ext    string
	Format string
	MIME   string
}{
	"image": {
		{"jpg", "1920x1080", "image/jpeg"},
		{"jpeg", "1920x1080", "image/jpeg"},
		{"png", "1280x720", "image/png"},
		{"webp", "1024x768", "image/webp"},
		{"gif", "800x600", "image/gif"},
	},
	"video": {
		{"mp4", "Full HD", "video/mp4"},
		// {"mov", "4K", "video/quicktime"},
		// {"webm", "HD", "video/webm"},
		// {"avi", "HD", "video/x-msvideo"},
	},
	"audio": {
		{"mp3", "320kbps", "audio/mpeg"},
		// {"wav", "Lossless", "audio/wav"},
		// {"ogg", "192kbps", "audio/ogg"},
		// {"aac", "256kbps", "audio/aac"},
	},
}

func ensureDirs(base string) {
	for _, typ := range []string{"image", "video", "audio"} {
		path := filepath.Join(base, typ)
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf("Erreur création dossier %s: %v", typ, err)
		}
	}
}

func randomMediaType() string {
	//nolint:gosec // math/rand suffisant ici
	return []string{"image", "video", "audio"}[rand.Intn(3)]
}

func createFakeFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("erreur création fichier %s: %w", path, err)
	}
	defer f.Close()
	_, err = f.WriteString("dummy content")
	if err != nil {
		return fmt.Errorf("erreur écriture contenu dans %s: %w", path, err)
	}
	return nil
}

func seedMedia(db *gorm.DB, userID uint64, count int, basePath string) []models.Media {
	//nolint:gosec // math/rand suffisant ici
	rand.New(rand.NewSource(time.Now().UnixNano()))
	ensureDirs(basePath)
	mediaList := make([]models.Media, 0, count)

	for i := 0; i < count; i++ {
		typ := randomMediaType()
		//nolint:gosec // math/rand suffisant ici
		def := mediaTypes[typ][rand.Intn(len(mediaTypes[typ]))]

		id := uuid.New().String()
		filename := fmt.Sprintf("%s.%s", id, def.Ext)
		dir := filepath.Join(basePath, typ)
		fullPath := filepath.Join(dir, filename)

		loadFiles(def, fullPath, filename)

		media := models.Media{
			FileName: filename,
			FileExt:  def.Ext,
			FileType: def.MIME,
			FileSize: 1024 * 1024,
			Model: models.Model{
				UUID: uuid.NullUUID{
					UUID:  uuid.New(),
					Valid: true,
				},
				CreatedAt: time.Now().UTC(),
				UpdatedAt: time.Now().UTC(),
			},
			FilePath: fullPath,
			UserID:   userID,
		}

		db.Create(&media)
		mediaList = append(mediaList, media)
	}

	fmt.Printf("✅ %d fichiers média créés dans '%s/' et insérés en base\n", count, basePath)
	return mediaList
}

func loadFiles(def struct {
	Ext    string
	Format string
	MIME   string
}, fullPath string, filename string,
) {
	var err error
	switch def.Ext {
	case "jpg":
		err = copyDummyFile("test/fixtures/dummy.jpg", fullPath)
	case "jpeg":
		err = copyDummyFile("test/fixtures/dummy.jpg", fullPath)
	case "webp":
		err = copyDummyFile("test/fixtures/dummy.webp", fullPath)
	case "gif":
		err = copyDummyFile("test/fixtures/dummy.gif", fullPath)
	case "png":
		err = copyDummyFile("test/fixtures/dummy.png", fullPath)
	case "mp4":
		err = copyDummyFile("test/fixtures/dummy.mp4", fullPath)
	case "mp3":
		err = copyDummyFile("test/fixtures/dummy.mp3", fullPath)
	default:
		err = createFakeFile(fullPath)
	}
	if err != nil {
		log.Fatalf("Erreur création fichier %s: %v", filename, err)
	}
}

func seedData(db *gorm.DB) {
	users := []models.User{
		{Username: "tester", Password: hashPassword("12345"), Bio: "Je suis Tester", Avatar: "", Role: "user", Email: "tester@mail.mail", Verified: true, VerifiedExpires: time.Now().Add(50 * 30 * 24 * time.Hour)},
		{Username: "alice", Password: hashPassword("12345"), Bio: "Je suis Alice", Avatar: "", Role: "user", Email: "alice@mail.mail"},
		{Username: "bob", Password: hashPassword("12345"), Bio: "Je suis Bob", Avatar: "", Role: "user", Email: "bob@mail.mail"},
		{Username: "medz", Password: hashPassword("12345"), Bio: "Je suis Medz", Avatar: "", Role: "user", Email: "medz@mail.mail"},
		{Username: "medz2", Password: hashPassword("12345"), Bio: "Je suis Medz2", Avatar: "", Role: "user", Email: "medz2@mail.mail"},
	}

	for i := range users {
		var existing models.User
		if err := db.Where("username = ?", users[i].Username).First(&existing).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Fatalf("Erreur recherche utilisateur %s: %v", users[i].Username, err)
			}
			db.Create(&users[i])
		} else {
			users[i].ID = existing.ID
		}

		m := seedMedia(db, users[i].ID, 5, "uploads")
		for j := 1; j <= 5; j++ {
			//nolint:gosec // math/rand suffisant ici
			now := time.Now().Add(-time.Duration(rand.Intn(100)) * time.Hour)
			post := models.Post{
				Content: fmt.Sprintf("Post #%d by %s", j, users[i].Username),
				UserID:  users[i].ID,
				Medias: []models.Media{
					m[i],
				},
				Model: models.Model{
					CreatedAt: now,
					UpdatedAt: now,
				},
			}
			db.Create(&post)
		}
	}

	var posts []models.Post
	db.Find(&posts)
	for idx := range posts {
		var like models.Like
		if err := db.Where("user_id = ? AND post_id = ?", users[0].ID, posts[idx].ID).First(&like).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				log.Fatalf("Erreur recherche like pour post %d: %v", posts[idx].ID, err)
			}
			db.Create(&models.Like{PostID: posts[idx].ID, UserID: users[0].ID})
		}
	}

	seedComments(db, posts, users)
	seedMessage(db, users)
}

func seedComments(db *gorm.DB, posts []models.Post, users []models.User) {
	comments := []string{
		"Super post !",
		"J'adore ce que tu fais.",
		"Très intéressant.",
		"Merci pour le partage !",
		"👌🔥",
		"Je suis d'accord avec toi.",
		"Top !",
		"Tu peux en dire plus ?",
	}

	for idx := range posts {
		//nolint:gosec // math/rand suffisant ici
		n := rand.Intn(10) + 1 // entre 1 et 3 commentaires par post
		for i := 0; i < n; i++ {
			comment := models.Comment{
				//nolint:gosec // math/rand suffisant ici
				Content: comments[rand.Intn(len(comments))],
				UserModel: models.UserModel{
					//nolint:gosec // math/rand suffisant ici
					UserID: users[rand.Intn(len(users))].ID,
				},
				PostID: posts[idx].ID,
			}
			db.Create(&comment)
		}
	}
}

func seedMessage(db *gorm.DB, users []models.User) {
	messageContents := []string{
		"Salut, comment tu vas ?",
		"Tu as vu mon dernier post ?",
		"On se capte ce soir ?",
		"Merci pour ton like !",
		"Top ton dernier média 👍",
	}

	for i := 0; i < 10; i++ {
		//nolint:gosec // math/rand suffisant ici
		sender := users[rand.Intn(len(users))]
		//nolint:gosec // math/rand suffisant ici
		receiver := users[rand.Intn(len(users))]

		if sender.ID == receiver.ID {
			continue
		}

		msg := models.Message{
			SenderID:   sender.ID,
			ReceiverID: receiver.ID,
			//nolint:gosec // math/rand suffisant ici
			Content: messageContents[rand.Intn(len(messageContents))],
			//nolint:gosec // math/rand suffisant ici
			Read:   rand.Intn(2) == 0,
			UserID: sender.ID,
			Model: models.Model{
				//nolint:gosec // math/rand suffisant ici
				CreatedAt: time.Now().Add(-time.Duration(rand.Intn(100)) * time.Hour),
				UpdatedAt: time.Now(),
			},
		}
		db.Create(&msg)

		notif := models.Notification{
			UserID:  receiver.ID,
			Content: fmt.Sprintf("Nouveau message de %s", sender.Username),
			IsRead:  false,
			Type:    models.NotificationTypeMessage,
			Link:    fmt.Sprintf("/chat/%d", sender.ID),
			TypeID:  msg.ID,
			Model: models.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}
		db.Create(&notif)
	}
}

func hashPassword(pw string) string {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(pw), 14)
	return string(hashed)
}

func copyDummyFile(srcPath, destPath string) error {
	src, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("erreur ouverture fichier source: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("erreur création fichier destination: %w", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return fmt.Errorf("erreur copie fichier: %w", err)
	}
	return nil
}

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Erreur chargement configuration: %v", err)
	}
	conn, err := connector.NewDBConn(cfg.DB)
	if err != nil {
		panic(fmt.Errorf("failed to create database connection: %w", err))
	}

	seedData(conn.DB)
}
