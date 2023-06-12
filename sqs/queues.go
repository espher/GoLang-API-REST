package sqs

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

func GetQueueList() {
	// Crea una sesión de AWS
	sesssion, err := CreateAWSSession()
	if err != nil {
		fmt.Printf(err.Error())
	}

	svc := sqs.New(sesssion)

	// Obtiene la lista de colas existentes
	result, err := svc.ListQueues(nil)
	if err != nil {
		fmt.Printf(err.Error())
	}

	for i, url := range result.QueueUrls {
		fmt.Printf("%d: %s\n", i, *url)
	}
}

// SendMessage envía un mensaje a una cola de AWS SQS.
func SendMessage(queueURL, message string) error {
	// Crea una sesión de AWS
	sess, err := CreateAWSSession()
	if err != nil {
		return err
	}

	// Crea un cliente de AWS SQS utilizando la sesión
	svc := sqs.New(sess)

	// Envía el mensaje a la cola especificada por la URL
	_, err = svc.SendMessage(&sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    &queueURL,
	})
	if err != nil {
		return err
	}

	return nil
}

// SendMessages envía mensajes en lote a una cola de AWS SQS.
func SendMessagesBulk(queueURL string, messages []string) error {
	// Crea una sesión de AWS
	sess, err := CreateAWSSession()
	if err != nil {
		return err
	}

	// Crea un cliente de AWS SQS utilizando la sesión
	svc := sqs.New(sess)

	// Prepara los mensajes para enviar en lote
	entries := make([]*sqs.SendMessageBatchRequestEntry, len(messages))
	for i, msg := range messages {
		entries[i] = &sqs.SendMessageBatchRequestEntry{
			Id:          aws.String(string(i)), // Identificador único para cada mensaje
			MessageBody: aws.String(msg),
		}
	}

	// Divide los mensajes en lotes de máximo 10 mensajes por solicitud
	batchSize := 10
	for i := 0; i < len(entries); i += batchSize {
		end := i + batchSize
		if end > len(entries) {
			end = len(entries)
		}

		batch := entries[i:end]

		// Envía el lote de mensajes a la cola especificada por la URL
		_, err := svc.SendMessageBatch(&sqs.SendMessageBatchInput{
			Entries:  batch,
			QueueUrl: &queueURL,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

// ReceiveMessages lee mensajes de una cola de AWS SQS.
func ReceiveMessages(queueURL string, maxMessages int64) ([]string, error) {
	sess, err := CreateAWSSession()
	if err != nil {
		return nil, err
	}

	svc := sqs.New(sess)

	// Recibe los mensajes de la cola especificada por la URL
	result, err := svc.ReceiveMessage(&sqs.ReceiveMessageInput{
		QueueUrl:            &queueURL,
		MaxNumberOfMessages: &maxMessages,
	})
	if err != nil {
		return nil, err
	}

	// Extrae el contenido de los mensajes recibidos
	var messages []string
	for _, msg := range result.Messages {
		messages = append(messages, *msg.Body)
	}

	return messages, nil
}

// DeleteMessages borra los mensajes especificados de una cola de AWS SQS.
func DeleteMessages(queueURL string, receiptHandles []*string) error {
	// Crea una sesión de AWS
	sess, err := CreateAWSSession()
	if err != nil {
		return err
	}

	// Crea un cliente de AWS SQS utilizando la sesión
	svc := sqs.New(sess)

	// Borra los mensajes de la cola especificada por la URL y los receipt handles
	_, err = svc.DeleteMessageBatch(&sqs.DeleteMessageBatchInput{
		QueueUrl: &queueURL,
		Entries:  createDeleteMessageBatchEntries(receiptHandles),
	})
	if err != nil {
		return err
	}

	return nil
}

// createDeleteMessageBatchEntries crea los objetos DeleteMessageBatchRequestEntry para la operación DeleteMessageBatch.
func createDeleteMessageBatchEntries(receiptHandles []*string) []*sqs.DeleteMessageBatchRequestEntry {
	var entries []*sqs.DeleteMessageBatchRequestEntry

	for i, receiptHandle := range receiptHandles {
		entry := &sqs.DeleteMessageBatchRequestEntry{
			Id:            aws.String(fmt.Sprintf("entry-%d", i)),
			ReceiptHandle: receiptHandle,
		}
		entries = append(entries, entry)
	}

	return entries
}
