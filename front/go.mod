module github.com/zitrator/stash_telegram/front

go 1.18

replace github.com/zitrator/stash_telegram/stash => ../stash

replace github.com/zitrator/stash_telegram/stashrest => ../stashrest

require (
	github.com/zitrator/stash_telegram/stash v0.0.0-20220611161032-0f4b6a29f339
	github.com/zitrator/stash_telegram/stashrest v0.0.0-00010101000000-000000000000
)

require github.com/gorilla/mux v1.8.0 // indirect
