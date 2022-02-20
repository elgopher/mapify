// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package mapify

func acceptAllFields(string, Element) bool {
	return true
}

func noRename(_ string, e Element) string {
	return e.Name()
}

func interfaceValue(_ string, e Element) interface{} {
	return e.Interface()
}
