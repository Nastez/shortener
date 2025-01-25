package saver

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Nastez/shortener/internal/app/models"
)

type Producer struct {
	file *os.File // файл для записи
	// добавляем Writer в Producer
	encoder *json.Encoder
}

type Consumer struct {
	file    *os.File
	decoder *json.Decoder
}

var Events = []*models.Event{
	{
		UUID:        1,
		ShortURL:    "4rSPg8ap",
		OriginalUrl: "http://yandex.ru",
	},
	{
		UUID:        2,
		ShortURL:    "edVPg3ks",
		OriginalUrl: "http://ya.ru",
	},
	{
		UUID:        3,
		ShortURL:    "dG56Hqxm",
		OriginalUrl: "http://practicum.yandex.ru",
	},
}

func NewProducer(fileName string) (*Producer, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return &Producer{
		file:    file,
		encoder: json.NewEncoder(file),
	}, nil
}

func (p *Producer) WriteEvent(event *models.Event) error {
	return p.encoder.Encode(&event)
}

func (p *Producer) Close() error {
	return p.file.Close()
}

func NewConsumer(fileName string) (*Consumer, error) {
	file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		file:    file,
		decoder: json.NewDecoder(file),
	}, nil
}

func (c *Consumer) ReadEvent() (*models.Event, error) {
	event := &models.Event{}
	if err := c.decoder.Decode(&event); err != nil {
		return nil, err
	}

	return event, nil
}

func (c *Consumer) Close() error {
	return c.file.Close()
}

func SaveFile(fileName string) {
	Producer, err := NewProducer(fileName)
	if err != nil {
		errors.New("can't open file")
	}
	defer Producer.Close()

	Consumer, err := NewConsumer(fileName)
	if err != nil {
		errors.New("can't open file")
	}
	defer Consumer.Close()

	for _, event := range Events {
		if err := Producer.WriteEvent(event); err != nil {
			errors.New("can't write event")
		}

		readEvent, err := Consumer.ReadEvent()
		if err != nil {
			errors.New("can't read file")
		}

		fmt.Println(readEvent)
	}
}
