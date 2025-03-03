import torch
print(torch.cuda.is_available())  # Должно вернуть True, если GPU доступен
print(torch.cuda.get_device_name(0))  # Название вашей видеокарты
