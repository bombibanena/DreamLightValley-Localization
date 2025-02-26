<a id="readme-top"></a>

[![Contributors][contributors-shield]][contributors-url]
[![Forks][forks-shield]][forks-url]
[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]
[![MIT License][license-shield]][license-url]

<br />
<div align="center">
  <h3 align="center">DreamLight Valley Localization</h3>

  <p align="center">
    Локализация для DreamLight Valley
    <br />
    <br />
    <a href="https://github.com/bombibanena/DreamLightValley-Localization/issues/new?labels=bug">Report Bug</a>
    &middot;
    <a href="https://github.com/bombibanena/DreamLightValley-Localization/issues/new?labels=enhancement">Request Feature</a>
  </p>
</div>

<details>
  <summary>Table of Contents</summary>
  <ol>
    <li>
      <a href="#готовый-машинный-перевод-от-меня">Готовый машинный перевод от меня</a>
      <ul>
        <li><a href="#версии">Версии</a></li>
      </ul>
    </li>
    <li>
      <a href="#для-тех-кто-хочет-сделать-собственный-перевод">Для тех кто хочет сделать собственный перевод</a>
    </li>
    <li>
      <a href="#decodeencode">Decode/Encode</a>
      <ul>
        <li><a href="#использование">Использование</a></li>
        <li><a href="#опции">Опции</a></li>
        <li><a href="#режим-decode">Режим decode</a></li>
        <li><a href="#режим-encode">Режим encode</a></li>
        <li><a href="#форматы">Форматы</a></li>
        <li><a href="#примеры">Примеры</a></li>
      </ul>
    </li>
    <li><a href="#license">License</a></li>
    <li><a href="#поддержать-автора">Поддержать автора</a></li>
  </ol>
</details>

## Готовый машинный перевод от меня

1. Загружаем архив для нужной [версии](#версии) (`Locale.zip`)
2. Переходим в корневую папку игры, там где расположен `ddv.exe`
3. Переходим в папку `ddv_Data\StreamingAssets\Localization`. Тут находятся все переводы на разных языках
4. Копируем файл `LocDB_en-US.zip` в любое другое место на компьютере. На случай если что-то пойдет не так, можно будет вернуться к оригинальному файлу
5. Заменяем файл `LocDB_en-US.zip` на такой же из загруженного архива `Locale.zip`
6. Запускаем игру и в настройках выбираем английский язык

Если нужна возможность переключаться обратно на английский - оригинальный `LocDB_en-US.zip` не трогаем. А файл из архива переименовываем в любой файл, который есть в папке `ddv_Data\StreamingAssets\Localization`. Например, в `LocDB_de.zip`. Но тогда и настройках нужно выбирать соответствующий язык

### Версии

- [v1.14.1.990 + 2 DLC](https://github.com/bombibanena/DreamLightValley-Localization/releases/tag/v.1.0.0-v1.14.1.990%2B2DLC)


## Для тех кто хочет сделать собственный перевод

1. Загружаем архив `Util.zip` из [Конвертер v.1.0.0](https://github.com/bombibanena/DreamLightValley-Localization/releases/tag/v.1.0.0), разархивируем в любом месте на ПК
2. Копируем файл с переводами из `ddv_Data\StreamingAssets\Localization`, например, `LocDB_en-US.zip` рядом с `ddv_loc.exe`
3. Разархивируем `LocDB_en-US.zip` в `LocDB_en-US`
4. Изучаем [как работать с ddv_loc.exe](#decodeencode)
5. Декодируем файлы в нужный формат
6. Делаем перевод
7. Кодируем файлы обратно
8. Архивируем файлы из `LocDB_en-US` в `LocDB_en-US.zip`
9. Заменяем `LocDB_en-US.zip` в `ddv_Data\StreamingAssets\Localization` новым архивом с переводами


## Decode/Encode

### Использование

```cmd
.\ddv_loc.exe --mode {decode|encode} --format {json|csv|raw} --in in_folder --out out_folder
```

### Опции
<table>
	<thead>
		<tr>
			<th>Флаг</th>
			<th>Описание</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td><code>--mode</code></td>
			<td>Режим работы: <code>decode</code> или <code>encode</code>.</td>
		</tr>
		<tr>
			<td><code>--format</code></td>
			<td>Формат данных: <code>json</code>, <code>csv</code> или <code>raw</code>.</td>
		</tr>
		<tr>
			<td><code>--in</code></td>
			<td>Входная папка или файл.</td>
		</tr>
		<tr>
			<td><code>--out</code></td>
			<td>Выходная папка.</td>
		</tr>
		<tr>
			<td><code>--help</code></td>
			<td>Показать справку.</td>
		</tr>
		<tr>
			<td><code>--version</code></td>
			<td>Показать версию программы.</td>
		</tr>
	</tbody>
</table>

### Режим decode
<table>
	<thead>
		<tr>
			<th>Флаг</th>
			<th>Описание</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td><code>--in</code></td>
			<td>Путь до папки с <code>.locbin</code> файлами.</td>
		</tr>
		<tr>
			<td><code>--out</code></td>
			<td>Результат работы в виде <code>loc.json</code>, <code>loc.csv</code> или папка с <code>.txt</code> файлами.</td>
		</tr>
	</tbody>
</table>

### Режим encode
<table>
	<thead>
		<tr>
			<th>Флаг</th>
			<th>Описание</th>
		</tr>
	</thead>
	<tbody>
		<tr>
			<td><code>--in</code></td>
			<td>Путь до <code>json</code>, <code>csv</code> или <code>.txt</code> файлов.</td>
		</tr>
		<tr>
			<td><code>--out</code></td>
			<td>Результат работы в виде <code>.locbin</code> файлов.</td>
		</tr>
	</tbody>
</table>

### Форматы

<h3><code>json</code></h3>

```json
[
  {
    "location": "/folder/file.locbin",
    "dictionary": [
      {
        "key": "key_1",
        "loc": {
          "en": "Value",
          "ru": "Значение"
        }
      }
    ]
  }
]
```

<h3><code>csv</code></h3>

```csv
location,key,en,ru
/folder/file.locbin,key_1,Value,Значение
```

<h3><code>raw</code></h3>

Декодированные `.locbin` файлы в `.txt` формате, с сохранением исходной структуры папок

### Примеры

<h3><code>json</code></h3>


* Декодирование:

    ```cmd
    .\ddv_loc.exe --mode decode --format json --in in_folder --out out_folder
    ```
* Кодирование:

    ```cmd
    .\ddv_loc.exe --mode encode --format json --in in_folder/loc.json --out out_folder
    ```

<h3><code>csv</code></h3>

* Декодирование:

    ```cmd
    .\ddv_loc.exe --mode decode --format csv --in in_folder --out out_folder
    ```
* Кодирование:

    ```cmd
    .\ddv_loc.exe --mode encode --format csv --in in_folder/loc.csv --out 
    ```

<h3><code>raw</code></h3>

* Декодирование:

    ```cmd
    .\ddv_loc.exe --mode decode --format raw --in in_folder --out out_folder
    ```

* Кодирование:

    ```cmd
    .\ddv_loc.exe --mode encode --format raw --in in_folder --out out_folder
    ```

## License

Distributed under the Unlicense License. See `LICENSE.txt` for more information.


## Поддержать автора

- BTC: `14scpxWeW6Lr7TDqcysgBvqjFGao43jYfw`
- TON: `UQCWzhdWmX463WqdKLLvQqssDEEmK_WFpS_I6C8-j56VR_o6`
- USDT (TRC20): `TNbeJkE6mu7xWWadwycPUQGSj7Bw6KCBo7`

<p align="right">(<a href="#readme-top">Наверх</a>)</p>


[contributors-shield]: https://img.shields.io/github/contributors/bombibanena/DreamLightValley-Localization.svg?style=for-the-badge
[contributors-url]: https://github.com/bombibanena/DreamLightValley-Localization/graphs/contributors
[forks-shield]: https://img.shields.io/github/forks/bombibanena/DreamLightValley-Localization.svg?style=for-the-badge
[forks-url]: https://github.com/bombibanena/DreamLightValley-Localization/network/members
[stars-shield]: https://img.shields.io/github/stars/bombibanena/DreamLightValley-Localization.svg?style=for-the-badge
[stars-url]: https://github.com/bombibanena/DreamLightValley-Localization/stargazers
[issues-shield]: https://img.shields.io/github/issues/bombibanena/DreamLightValley-Localization.svg?style=for-the-badge
[issues-url]: https://github.com/bombibanena/DreamLightValley-Localization/issues
[license-shield]: https://img.shields.io/github/license/bombibanena/DreamLightValley-Localization.svg?style=for-the-badge
[license-url]: https://github.com/bombibanena/DreamLightValley-Localization/blob/master/LICENSE.txt

