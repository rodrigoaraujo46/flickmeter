package user

import (
	"fmt"
	"math/rand"
	"time"
)

type User struct {
	Id        uint      `json:"id"`
	Email     string    `json:"-"`
	Username  string    `json:"username"`
	AvatarURL string    `json:"avatar_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func New(email, username, avatar_url string) *User {
	return &User{
		Email:     email,
		Username:  username,
		AvatarURL: avatar_url,
	}
}

func (u *User) SetRandomUsername() {
	adjectives := []string{
		"Swift", "Fast", "Cool", "Clever", "Bright", "Bold", "Lucky", "Chill",
		"Creative", "Brave", "Happy", "Fierce", "Gentle", "Mighty", "Sly", "Wise",
		"Nimble", "Sharp", "Witty", "Energetic", "Quick", "Silent", "Jolly",
		"Daring", "Fearless", "Radiant", "Glorious", "Honest", "Lively", "Curious",
		"Majestic", "Valiant", "Vivid", "Playful", "Serene", "Heroic", "Gracious",
		"Cheerful", "Cunning", "Gallant", "Luminous", "Noble", "Optimistic",
		"Persistent", "Rebellious", "Vibrant", "Whimsical", "Zesty", "Zealous",
		"Agile", "Alert", "Amused", "Bouncy", "Brilliant", "Calm", "Charming",
		"Dazzling", "Delightful", "Determined", "Eager", "Fabulous", "Faithful",
		"Fancy", "Fearless", "Flawless", "Friendly", "Funny", "Gentle", "Gleaming",
		"Glorious", "Graceful", "Happy", "Honest", "Inspiring", "Jovial", "Kind",
		"Lucky", "Magical", "Mighty", "Optimistic", "Peaceful", "Playful", "Proud",
		"Radiant", "Respectful", "Sassy", "Sincere", "Smart", "Snappy", "Sociable",
		"Strong", "Stylish", "Sunny", "Talented", "Thankful", "Unique", "Upbeat",
		"Valiant", "Vivid", "Wise", "Witty", "Youthful", "Zealous",
	}

	animals := []string{
		"Lion", "Tiger", "Falcon", "Bear", "Wolf", "Panda", "Shark", "Eagle",
		"Fox", "Hawk", "Panther", "Otter", "Cheetah", "Jaguar", "Dolphin", "Lynx",
		"Raven", "Stallion", "Buffalo", "Cobra", "Moose", "Badger", "Cougar",
		"Elephant", "Hippo", "Gorilla", "Kangaroo", "Leopard", "Raccoon",
		"Seahorse", "Swan", "TigerShark", "Vulture", "Walrus", "Zebra", "Alligator",
		"Bison", "Camel", "Dragon", "Fennec", "Gazelle", "Heron", "Iguana", "Jackal",
		"Koala", "Marmot", "Narwhal", "Ocelot", "Penguin", "Quokka", "Rattlesnake",
		"Salamander", "Toucan", "Urial", "Viper", "Wombat", "Xerus", "Yak", "Zorilla",
		"Albatross", "Barracuda", "Caribou", "Donkey", "Emu", "Ferret", "Giraffe",
		"Hedgehog", "Impala", "Jay", "Kiwi", "Lemur", "Mole", "Numbat", "Ostrich",
		"Porcupine", "Quail", "Rabbit", "Salmon", "Tapir", "Urchin", "Vicuña",
		"Walabi", "Xantus", "Yellowtail", "Zebu", "Armadillo", "Beaver", "Catfish",
		"Dugong", "Eland", "Flamingo", "Goose", "Hummingbird", "Ibex", "Jellyfish",
		"Kudu", "Lobster", "Manatee", "Newt", "Octopus", "Platypus", "Quetzal",
		"Reindeer", "Seagull", "Tarantula", "Urutu", "VultureBat", "Warthog",
	}

	adjIndex := rand.Intn(len(adjectives))
	animalIndex := rand.Intn(len(animals))
	num := rand.Intn(999_999)

	u.Username = fmt.Sprintf("%s%s%d", adjectives[adjIndex], animals[animalIndex], num)
}
