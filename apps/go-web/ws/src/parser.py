
def do_parse():
  methods = set()
  with open('live.douyin.com-1750573630444.log', encoding='utf-8') as fp:
    for line in fp.readlines():
      line = line.strip()
      if line == '':
        continue
      segments = line.split(' -> ')
      method = segments[1].replace(', t', '')
      body = segments[2]
      methods.add(method)
      with open(f'{method}.json', 'w', encoding='utf-8') as fp:
        fp.write(body)
    pass
  ms = '[]string{'
  for method in methods:
    ms += f'"{method}",'
  ms += '}'
  print(f'all methods size: {len(methods)} -> {methods} -> {ms}')


if __name__ == '__main__':
  do_parse()
