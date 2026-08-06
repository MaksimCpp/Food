#!/usr/bin/env python3
import os
import sys
from pathlib import Path

def create_service_structure():
    """
    Создает структуру папок и пустых файлов для каждого сервиса в папке services.
    """
    services_dir = "services"
    
    # Проверяем существование папки services
    if not os.path.exists(services_dir):
        print(f"Ошибка: Папка '{services_dir}' не найдена!", file=sys.stderr)
        sys.exit(1)
    
    if not os.path.isdir(services_dir):
        print(f"Ошибка: '{services_dir}' не является директорией!", file=sys.stderr)
        sys.exit(1)
    
    # Получаем список всех подпапок в services
    try:
        services = [d for d in os.listdir(services_dir) 
                   if os.path.isdir(os.path.join(services_dir, d))]
    except PermissionError:
        print(f"Ошибка: Нет прав на чтение папки '{services_dir}'", file=sys.stderr)
        sys.exit(1)
    
    if not services:
        print(f"Предупреждение: В папке '{services_dir}' нет поддиректорий")
        return
    
    print(f"Найдены сервисы: {', '.join(services)}")
    print("-" * 60)
    
    # Структура папок и файлов (только названия)
    structure = {
        "cmd": {
            "server": ["main.go"]
        },
        "internal": {
            "config": ["config.go"],
            "delivery": ["delivery.go"],
            "grpc": ["server.go"],
            "domain": ["domain.go"],
            "repository": ["repository.go"],
            "postgres": ["postgres.go"],
            "redis": ["redis.go"],
            "infrastructure": ["infrastructure.go"],
            "kafka": ["producer.go", "consumer.go"],
            "usecase": ["usecase.go"]
        },
        "proto": ["service.proto"]
    }
    
    # Файлы в корне сервиса
    root_files = ["Dockerfile"]
    
    total_created = 0
    
    # Проходим по каждому сервису
    for service in services:
        service_path = Path(services_dir) / service
        print(f"\n📁 Создаем структуру для {service}:")
        
        created = create_structure(service_path, structure, root_files, service)
        total_created += created
        
        # Создаем go.mod
        create_go_mod(service_path, service)
        
        print(f"   ✅ Создано {created} файлов/папок")
    
    print("-" * 60)
    print(f"✅ Всего создано {total_created} файлов/папок во всех сервисах!")

def create_structure(base_path, structure, root_files, service_name):
    """
    Рекурсивно создает структуру папок и пустых файлов.
    """
    created_count = 0
    
    # Создаем файлы в корне
    for filename in root_files:
        file_path = base_path / filename
        if not file_path.exists():
            try:
                file_path.touch()
                created_count += 1
                print(f"   📄 {file_path.relative_to(base_path.parent)}")
            except Exception as e:
                print(f"   ❌ Ошибка при создании {file_path}: {str(e)}")
        else:
            print(f"   ⚠️  {file_path.relative_to(base_path.parent)} уже существует, пропускаем")
    
    # Создаем папки и файлы по структуре
    for name, content in structure.items():
        current_path = base_path / name
        
        if isinstance(content, dict):
            # Это папка с вложенными папками
            try:
                current_path.mkdir(exist_ok=True, parents=True)
                created_count += 1
                print(f"   📁 {current_path.relative_to(base_path.parent)}")
                
                # Создаем вложенную структуру
                for sub_name, files in content.items():
                    sub_path = current_path / sub_name
                    sub_path.mkdir(exist_ok=True, parents=True)
                    created_count += 1
                    print(f"   📁 {sub_path.relative_to(base_path.parent)}")
                    
                    # Создаем файлы в вложенной папке
                    for filename in files:
                        file_path = sub_path / filename
                        if not file_path.exists():
                            file_path.touch()
                            created_count += 1
                            print(f"   📄 {file_path.relative_to(base_path.parent)}")
                        else:
                            print(f"   ⚠️  {file_path.relative_to(base_path.parent)} уже существует, пропускаем")
                            
            except PermissionError:
                print(f"   ❌ Нет прав на создание папки {current_path}")
        else:
            # Это список файлов в папке
            try:
                current_path.mkdir(exist_ok=True, parents=True)
                created_count += 1
                print(f"   📁 {current_path.relative_to(base_path.parent)}")
                
                for filename in content:
                    file_path = current_path / filename
                    if not file_path.exists():
                        file_path.touch()
                        created_count += 1
                        print(f"   📄 {file_path.relative_to(base_path.parent)}")
                    else:
                        print(f"   ⚠️  {file_path.relative_to(base_path.parent)} уже существует, пропускаем")
                        
            except PermissionError:
                print(f"   ❌ Нет прав на создание папки {current_path}")
    
    return created_count

def create_go_mod(service_path, service_name):
    """
    Создает go.mod для сервиса.
    """
    go_mod_path = service_path / "go.mod"
    
    # Проверяем, существует ли уже go.mod
    if go_mod_path.exists():
        print(f"   ℹ️  go.mod уже существует для {service_name}")
        return
    
    try:
        import subprocess
        
        module_name = f"github.com/MaksimCpp/{service_name}"
        
        # Выполняем команду go mod init
        result = subprocess.run(
            ["go", "mod", "init", module_name],
            cwd=str(service_path),
            capture_output=True,
            text=True,
            check=True
        )
        print(f"   ✅ go.mod создан для {service_name} (модуль: {module_name})")
        
    except subprocess.CalledProcessError as e:
        print(f"   ❌ Ошибка при создании go.mod для {service_name}: {e.stderr.strip()}")
    except FileNotFoundError:
        print(f"   ❌ Команда 'go' не найдена. Убедитесь, что Go установлен.")
        # Создаем пустой go.mod вручную
        try:
            with open(go_mod_path, 'w', encoding='utf-8') as f:
                f.write(f"module github.com/MaksimCpp/{service_name}\n\ngo 1.21\n")
            print(f"   ⚠️  Создан базовый go.mod для {service_name}")
        except:
            pass
    except Exception as e:
        print(f"   ❌ Неожиданная ошибка при создании go.mod: {str(e)}")

if __name__ == "__main__":
    create_service_structure()