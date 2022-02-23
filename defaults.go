// (c) 2022 Jacek Olszak
// This code is licensed under MIT license (see LICENSE for details)

package mapify

func acceptAllFields(string, Element) (bool, error) {
	return true, nil
}

func noRename(_ string, e Element) (string, error) {
	return e.Name(), nil
}

func interfaceValue(_ string, e Element) (interface{}, error) {
	return e.Interface(), nil
}
