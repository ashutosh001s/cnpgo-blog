import os

def rename_tmpl_to_html(root_dir):
    for dirpath, dirnames, filenames in os.walk(root_dir):
        for filename in filenames:
            if filename.endswith('.tmpl'):
                tmpl_path = os.path.join(dirpath, filename)
                html_path = os.path.join(dirpath, filename[:-5] + '.html')
                os.rename(tmpl_path, html_path)
                print(f"Renamed: {tmpl_path} -> {html_path}")

if __name__ == "__main__":
    target_directory = input("Enter the path to the directory: ").strip()
    if os.path.isdir(target_directory):
        rename_tmpl_to_html(target_directory)
    else:
        print("Invalid directory path.")
