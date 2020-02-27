import sys
word=sys.argv[1]
max=0
for x in range(0,len(word)):
    characters=[]
    for y in range(x,len(word)):
        if word[y] in characters:
            if max< len(characters):
                max=len(characters)
            break
        else:
            characters.append(word[y])
print(max)                