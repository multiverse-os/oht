package database

import ()

// Translation is a struct used to save translations into databae
type Translation struct {
	Locale string `json:`
	Key    string `json:`
	Value  string `json:`
}

// LoadTranslations load translations from DB backend
//func (backend *Backend) LoadTranslations() (translations []*i18n.Translation) {
//	backend.DB.Find(&translations)
//	return translations
//}
//
//// SaveTranslation save translation into DB backend
//func (backend *Backend) SaveTranslation(t *i18n.Translation) error {
//	return backend.DB.Where(Translation{Key: t.Key, Locale: t.Locale}).
//		Assign(Translation{Value: t.Value}).
//		FirstOrCreate(&Translation{}).Error
//}
//
//// DeleteTranslation delete translation into DB backend
//func (backend *Backend) DeleteTranslation(t *i18n.Translation) error {
//	return backend.DB.Where(Translation{Key: t.Key, Locale: t.Locale}).Delete(&Translation{}).Error
//}
