package services

import (
	"context"
	"document-service/models"
	"document-service/store"
	"errors"
)

var Documents map[uint64]models.Document

// DB — репозиторий документов/пользователей, инициализируется в main.
var DB *store.Store

// SaveDocument сохраняет документ и связи с пользователями, которым он выдан.
func SaveDocument(ctx context.Context, doc models.Document) (models.Document, error) {
	_, ok := Session[doc.Token]
	if !ok {
		return models.Document{}, errors.New("no token in document")
	}
	document, err := DB.SaveDocument(ctx, doc)
	if err != nil {
		return document, err
	}
	for _, item := range doc.Grant {
		// Выбрать id пользователя по логину
		userID, err := DB.GetUserIDByLogin(ctx, item)
		if err != nil {
			return document, err
		}
		// Записать в userFileLink id пользователя и document.id
		err = DB.SaveUserFileLink(ctx, models.UsersFileLink{
			IdUser:     userID,
			IdDocument: document.ID,
		})
		if err != nil {
			return document, err
		}
	}
	// Сохранение на диск
	if err := document.SaveFile(); err != nil {
		return document, err
	}
	return document, nil
}

// GetListDocuments возвращает список документов по фильтрам
func GetListDocuments(ctx context.Context, limit, offset uint64) ([]models.Document, error) {
	result, err := DB.GetDocumentList(ctx, limit, offset)
	if err != nil {
		return result, err
	}
	for i, document := range result {
		for _, grant := range document.UsersFileLinkList {
			result[i].Grant = append(result[i].Grant, grant.Login)
		}
	}
	return result, nil
}

// GetDocument возвращает документ по ID
func GetDocument(ctx context.Context, id uint64) (models.Document, error) {
	doc, ok := Documents[id]
	if ok {
		return doc, nil
	}
	document, err := DB.GetDocument(ctx, id)
	if err != nil {
		return document, err
	}
	for _, grant := range document.UsersFileLinkList {
		document.Grant = append(document.Grant, grant.Login)
	}
	Documents[id] = document
	return document, nil
}

// DeleteDocument удаляет документ по ID
func DeleteDocument(ctx context.Context, id uint64) error {
	if err := DB.DeleteDocument(ctx, id); err != nil {
		return err
	}
	delete(Documents, id)
	return nil
}
