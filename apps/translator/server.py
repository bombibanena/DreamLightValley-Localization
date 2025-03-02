from flask import Flask, request, jsonify
from translator import translate, calcTokens

app = Flask(__name__)

@app.route('/translate', methods=['POST'])
def makeTranslateReq():
    data = request.json
    text = data.get('text', '')

    data = translate(text)

    return jsonify({'data': data})

@app.route('/tokens/calc', methods=['POST'])
def calcTokensReq():
    data = request.json
    text = data.get('text', '')

    data = calcTokens(text)

    return jsonify({'tokens': data})

# Запуск сервера
if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
