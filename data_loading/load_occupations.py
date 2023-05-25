from datetime import datetime

import requests
from pymongo import MongoClient
from tqdm import tqdm

# Define ONET API endpoint and parameters
onet_api_url = "https://services.onetcenter.org/ws/online/occupations/"
onet_api_params = {
    "start": 1,
    "end": 1016,
}

# Read in API keys and MongoDB Conenction URL with a .env file
with open('.env', 'r') as file:
    base64_credentials = file.readline().split('=')[1].replace("\"", "").strip()
    mongo_client_uri = file.readline().split('=')[1].replace("\"", "").strip()
    mongo_db_name = file.readline().split('=')[1].replace("\"", "").strip()

# Define MongoDB connection and collection
mongo_occupations_collection_name = "occupation_data"

# Create MongoDB Client and set up collections
mongo_client = MongoClient(mongo_client_uri)
mongo_db = mongo_client[mongo_db_name]
mongo_occupations_collection = mongo_db[mongo_occupations_collection_name]

# Define bearer token
headers = {
    "Authorization": f"Basic {base64_credentials}",
    'User-Agent': 'python-OnetWebService/1.00 (bot)',
    'Accept': 'application/json',
}

# Gets all occupation codes from ONET and stores them in two files for easy access
def get_all_occupation_codes():
    response = requests.get(onet_api_url, params=onet_api_params, headers=headers)
    if response.status_code == 200:
        data = response.json()

        occupation_codes = []
        occupation_urls = []
        for occupation in data["occupation"]:
            occupation_codes.append(occupation["code"])
            occupation_urls.append(occupation["href"])
        
        write_lines_to_file('data_loading/occupation_codes.txt', occupation_codes)
        # write_lines_to_file('data_loading/occupation_urls.txt', occupation_urls)

        print("Occupation data written to files")
    else:
        print("Failed to retrieve occupation data from ONET API.")

def write_lines_to_file(filename, lines):
    with open(filename, 'w') as file:
        file.writelines(line + "\n" for line in lines)

# Gets an occupation with the given occupation code from ONET and stores it in the raw_occupations_collection
def get_occupation(occupation_code):
    response = requests.get(onet_api_url + occupation_code + '/summary', params={'display': 'long'}, headers=headers)
    if response.status_code == 200:
        data = response.json()
        occupation_data = filter_occupation_data(data)
        workstyle_data = get_workstyles(occupation_code)
        occupation_data['work_styles'] = workstyle_data
        occupation_data['date_updated'] = datetime.utcnow()
        mongo_occupations_collection.replace_one({'code': occupation_code}, occupation_data, upsert=True)
    else:
        print(f"Failed to retrieve occupation data for {occupation_code} from ONET API.")
        print(f'Reason: {response.status_code}: {response.reason}')

def get_workstyles(occupation_code):
    response = requests.get(onet_api_url + occupation_code + '/details/work_styles', params={'display': 'long'}, headers=headers)
    if response.status_code == 200:
        data = response.json()
        workstyle_data = filter_workstyle_data(data)
        return workstyle_data
    else:
        return []

# Gets all occupations from ONET and stores it in the raw_occupations_collection
def get_all_occupation_data():
    with open('data_loading/occupation_codes.txt', 'r') as file:
        occupation_codes = file.readlines()
        for occupation_code in tqdm(occupation_codes):
            occupation_code = occupation_code.strip()
            if occupation_code != '': # and not occupation_exists_in_db(occupation_code):
                get_occupation(occupation_code)

# Converts a raw_occupation_data to occupation data that we will actually use and need
def filter_occupation_data(raw_occupation_data):
    occupation_data = {}
    occupation_data['code'] = raw_occupation_data['code']
    occupation_data['name'] = raw_occupation_data['occupation']['title']
    occupation_data['description'] = raw_occupation_data['occupation']['description']
    occupation_data['technology_skills'] = filter_tech_skills(raw_occupation_data.get('technology_skills', {'category': []}))
    return occupation_data

# Converts a raw technology skills to tech skills data
def filter_tech_skills(raw_tech_skills):   
    tech_skills = []
    for skill in raw_tech_skills['category']:
        data = {}
        data['id'] = skill['title']['id']
        data['name'] = skill['title']['name']
        data['examples'] = []
        for example in skill['example']:
            data['examples'].append(example['name'])
        tech_skills.append(data)
    return tech_skills

# Converts a raw_workstyle_data to workstyle data that we will actually use and need
def filter_workstyle_data(raw_workstyle_data):   
    workstyle_data = []
    for workstyle in raw_workstyle_data['element']:
        data = {}
        data['id'] = workstyle['id']
        data['name'] = workstyle['name']
        data['description'] = workstyle['description']
        data['scale'] = workstyle['score']['scale']
        data['value'] = workstyle['score']['value']
        workstyle_data.append(data)
    return workstyle_data

if __name__ == '__main__':
    get_all_occupation_data()
    print("All occupation data stored in MongoDB!")
