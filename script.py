#!/usr/bin/env python3
import os
import sys
from pathlib import Path

def delete_go_files():
    """
    Удаляет все .go файлы во всем проекте.
    """
    project_root = Path.cwd()
    
    print(f"📂 Сканируем проект: {project_root}")
    print("-" * 60)
    
    # Список для хранения найденных .go файлов
    go_files = []
    
    # Рекурсивно ищем все .go файлы
    try:
        for file_path in project_root.rglob("*.go"):
            go_files.append(file_path)
    except PermissionError:
        print("❌ Ошибка: Нет прав на чтение некоторых директорий", file=sys.stderr)
        sys.exit(1)
    
    if not go_files:
        print("✅ .go файлы не найдены")
        return
    
    print(f"Найдено {len(go_files)} .go файлов:")
    
    # Выводим список найденных файлов
    for file_path in go_files:
        relative_path = file_path.relative_to(project_root)
        print(f"   📄 {relative_path}")
    
    print("-" * 60)
    
    # Запрашиваем подтверждение
    response = input(f"Удалить все {len(go_files)} .go файлов? (y/N): ").strip().lower()
    
    if response not in ['y', 'yes']:
        print("❌ Операция отменена")
        return
    
    # Удаляем файлы
    deleted_count = 0
    error_count = 0
    
    for file_path in go_files:
        try:
            file_path.unlink()
            deleted_count += 1
            print(f"✅ Удален: {file_path.relative_to(project_root)}")
        except Exception as e:
            error_count += 1
            print(f"❌ Ошибка при удалении {file_path.relative_to(project_root)}: {str(e)}")
    
    print("-" * 60)
    print(f"✅ Удалено файлов: {deleted_count}")
    if error_count > 0:
        print(f"⚠️  Ошибок при удалении: {error_count}")

if __name__ == "__main__":
    delete_go_files()