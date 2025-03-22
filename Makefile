up:
	docker run -d \
		-e API_TOKEN=${API_TOKEN} \
		-e AUDIO_URL=${AUDIO_URL} \
		-e AUDIOFILES_PATH=${AUDIOFILES_PATH} \
		-v ./audiofiles.txt:/app/audiofiles.txt \
		tempestmon/kingpin_bot
