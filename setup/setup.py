import requests
import json
import argparse
import os

script_dir = os.path.dirname(__file__)

def metadata(name):
    filename = os.path.join(script_dir, f'metadata/{name}.json')
    with open(filename) as file:
        return json.load(file)

def post(args, data):
    url = f'{args.endpoint}/v1/query'
    headers = {'x-hasura-admin-secret': args.admin_secret}
    return requests.post(url, data=json.dumps(data), headers=headers)

def execute(args, name, action, prepare=None):
    data = metadata(action)
    if prepare is not None:
        data = prepare(data)
    out = post(args, data)
    if out.status_code == 400:
        code = out.json()['code']
        if data['type'] == 'create_action':
            data['type'] = 'update_action'
            out = post(args, data)
        if data['type'] == 'create_action_permission' and code == 'already-exists':
            print(name, code)
            return
    if out.status_code != 200:
        raise Exception(f'{out.status_code} {out.text}')
    print(name, out.json())

def prepare_action_permission(action, role):
    def prepare(json):
        json['args']['action'] = action
        json['args']['role'] = role
        return json
    return lambda json: prepare(json)

def setup(args):
    action = 'custom_types'
    execute(args, action, action)

    actions = ['authenticate', 'crypt', 'whoami']
    roles = ['user', 'manager', 'anonymous']
    for action in actions:
        execute(args, action, action)
        for role in roles:
            prepare = prepare_action_permission(action, role)
            execute(args, f'permission {action} {role}', 'create_action_permission', prepare)

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--endpoint", help="http(s) endpoint for Hasura GraphQL Engine", required=True)
    parser.add_argument("--admin-secret", help="admin secret for Hasura GraphQL Engine", required=True)
    args = parser.parse_args()
    setup(args)

if __name__ == "__main__":
    main()
