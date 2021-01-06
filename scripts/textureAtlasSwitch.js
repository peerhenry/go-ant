const chars = '?!"#$%&\'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~'
for(let n = 0; n<94; n++) {
  console.log(`  case '${chars[n]}':`)
  console.log(`    return ${n}`)
}