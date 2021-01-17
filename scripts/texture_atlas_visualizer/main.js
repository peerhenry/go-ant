(function() {
  const canvas = document.getElementById('canvas')
  const ctx = canvas.getContext('2d')
  drawImage(ctx)

  function drawImage(ctx) {
    window.addEventListener('load', function () {
      console.log("It's loaded!")
      const image = document.getElementById('image')
      ctx.drawImage(image, 0, 0 );
      drawGrid(ctx, canvas.height/16)
      drawCoords(ctx, canvas.height/16)
    })
  }

  function drawVerticalLine(ctx, x) {
    ctx.beginPath();
    ctx.strokeStyle = "#FFFFFF";
    ctx.moveTo(x,0);
    ctx.lineTo(x,canvas.height);
    ctx.stroke();
  }
  
  function drawHorizontalLine(ctx, y) {
    ctx.beginPath();
    ctx.strokeStyle = "#FFFFFF";
    ctx.moveTo(0,y);
    ctx.lineTo(canvas.width,y);
    ctx.stroke();
  }
  
  function drawGrid(ctx, cellSize) {
    for(let i = 1; i<16; i++) {
      drawVerticalLine(ctx, i*cellSize)
    }
    for(let j = 1; j<16; j++) {
      drawHorizontalLine(ctx, j*cellSize)
    }
  }

  function drawCoords(ctx, cellSize) {
    ctx.font = '10px arial';
    ctx.fillStyle = 'white'
    for(let i = 0; i<16; i++) {
      for(let j = 0; j<16; j++) {
        ctx.fillText(`${i}, ${j}`, (i+0.1)*cellSize, (j+0.3)*cellSize);
      }
    }
  }
})()
