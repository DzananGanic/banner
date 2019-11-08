package domain

// An example to show how banner API interface is used:
// In this example we are using postgres repository implementation
// and basic banner displaying algorithm from platform/displayer/basic.go
// For caching active banner, we use redis implementation as an example

// As code is decoupled, we can switch and use different repository implementations

/*

db := postgres.NewBannerDB()
aProvider := redis.NewBannerActiveProvider()

// creation of banner API
b := banner.New(
	db,
	displayer.NewBasic(
		db,
		aProvider,
		ip.Internal,
	),
)

// creating a new banner
resp, err := b.Create(
	context.Background(),
	banner.CreateReq{
		Name: "sample banner",
		ScheduledDisplayingAt: time.Now(),
		ExpiresAt: time.Now().AddDate(time.Hour*24),
	},
)
resp.ID

// updating existing banner
err := b.Update(
	context.Background(),
	banner.UpdateReq{
		ID: 1,
		...
	},
)

// calling .Display returns available and active domain banner
activeBanner, err := b.Display(context.Background())

*/
