from transformers import pipeline, MarianTokenizer
import re
import torch

global pipe
pipe = None

def translate(text):
    pipe = getPipe()

    s0_pattern = re.compile(r'\(S_0\)')
    ph_pattern = re.compile(r'\(PH_\d+\)')  # Для поиска (PH_*)
    q_pattern = re.compile(r'\(Q_\d+\)')
    
    # Регулярное выражение для проверки знаков препинания
    punctuation_pattern = re.compile(r'^[\W_]+$')  # Знаки препинания и символы, не являющиеся буквами или цифрами
    
    # Регулярное выражение для проверки возвращаемых значений из pipe
    pipe_return_pattern = re.compile(r'@[pPРр]\d+')  # Шаблоны @p<число>, @P<число>, @Р<число>, @р<число>
    
    # Разбиваем текст по (S_0)
    parts = s0_pattern.split(text)
    
    result = []
    
    for part in parts:
        # Создаем локальный список замен для (PH_*) и счётчик
        ph_replacements = {}
        ph_counter = 0
        
        # Функция для замены (PH_*) на ключи
        def replace_ph(match):
            nonlocal ph_counter
            key = f"@ph{ph_counter}"  # Уникальный ключ в формате @ph<число>
            ph_replacements[key] = match.group(0)  # Сохраняем оригинальное значение
            ph_counter += 1
            return key
        
        # Заменяем все (PH_*) на ключи в текущей части
        part = ph_pattern.sub(replace_ph, part)
        
        # Если часть содержит (Q_*), разбиваем на sub_parts
        if q_pattern.search(part):
            sub_parts = []
            last_end = 0
            for match in re.finditer(r'\(Q_\d+\)', part):
                start, end = match.span()
                if last_end < start:
                    sub_parts.append(part[last_end:start])
                sub_parts.append(match.group())
                last_end = end
            if last_end < len(part):
                sub_parts.append(part[last_end:])
            
            # Обрабатываем каждую sub_part
            translated_sub_parts = []
            for sub_part in sub_parts:
                if q_pattern.match(sub_part):
                    # Если это (Q_*), оставляем без изменений
                    translated_sub_parts.append(sub_part)
                elif punctuation_pattern.match(sub_part):
                    # Если это знак препинания, оставляем без изменений
                    translated_sub_parts.append(sub_part)
                else:
                    # Иначе переводим
                    translated_sub_parts.append(pipe(sub_part)[0]['translation_text'])
            
            # Собираем обратно
            part = ''.join(translated_sub_parts)
        
        # Если часть не содержит (Q_*), просто переводим
        else:
            part = pipe(part)[0]['translation_text']
        
        # Восстанавливаем (PH_*) из ключей, избегая замены возвращаемых значений из pipe
        for key, value in ph_replacements.items():
            # Заменяем только ключи, которые не являются возвращаемыми значениями из pipe
            if not pipe_return_pattern.search(part):
                part = part.replace(key, value)
        
        result.append(part)
    
    # Собираем все части обратно в одну строку с разделителями (S_0)
    return '(S_0)'.join(result)
      

    # q_pattern = re.compile(r'\(Q_0\)')
    # text = q_pattern.sub(r'"', text)

    # q1_pattern = re.compile(r'\(Q_1\)')
    # text = q1_pattern.sub(r' @q1 ', text)

    # separator = "(S_0)"
    # parts = text.split(separator)

    # placeholder_pattern = re.compile(r'\(PH_(\d+)\)')
    # parts = [placeholder_pattern.sub(r' @\1 ', part) for part in parts]

    # processed_parts = [(pipe(part))[0]['translation_text'] for part in parts]

    # res = separator.join(processed_parts)

    # res = re.sub(r'@(\d+)', r'(PH_\1)', res, flags=re.IGNORECASE)
    # # res = re.sub(r'р(\d+)', r'(PH_\1)', res, flags=re.IGNORECASE)
    # # res = re.sub(r'п(\d+)', r'(PH_\1)', res, flags=re.IGNORECASE)
    # res = re.sub(r'@q1', r'(Q_1)', res)
    # res = re.sub(r'"', r'(Q_0)', res)

    # return res

def calcTokens(text):
    tokenizer = MarianTokenizer.from_pretrained("Helsinki-NLP/opus-mt-en-ru")
    tokens = tokenizer.tokenize(text)
    return len(tokens)

def getPipe():
    global pipe 
    if pipe is None:
        model = "Helsinki-NLP/opus-mt-en-ru"
        # model="facebook/nllb-200-distilled-600M"
        device = 0 if torch.cuda.is_available() else -1 # (GPU: 0, CPU: -1)
        pipe = pipeline("translation", model=model, max_length=512, num_beams=5, device=device)
    
    return pipe
