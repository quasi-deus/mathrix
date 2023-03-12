var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}
var coll = document.getElementsByClassName("collapsible");
var i;

for (i = 0; i < coll.length; i++) {
	  coll[i].addEventListener("click", function() {
	      this.classList.toggle("active");
	      var content = this.nextElementSibling;
	      if (content.style.display === "block") {
		    content.style.display = "none";
		  } else {
	      content.style.display = "block";
	      }
	  });
}



let canvas = document.querySelector("canvas");

let ctx = canvas.getContext("2d");
let width = canvas.width = window.innerWidth;
let height = canvas.height = window.innerHeight;
//let str = "A+jk js:2 @dfs 17 tr YY ufds M5r P87 #18 $!& ^dfs $Ew er JH # $ * . (! ;) ,: :";
let str = "0 1";
let matrix = str.split("");
let font = 16;
let col = width / font;
let arr = [];

for(let i = 0; i < col; i++) {
	    arr[i] = 1;
}

const draw = () => {
	    ctx.fillStyle = "rgba(0,0,0,0.05)";
	    ctx.fillRect(0, 0, width, height);
	    ctx.fillStyle = "#00FF00";
	    ctx.font = `${font}px system-ui`;

	    for(let i = 0; i < arr.length; i++) {
		            let txt = matrix[Math.floor(Math.random() * matrix.length)];
		            ctx.fillText(txt, i * font, arr[i] * font);

		            if(arr[i] * font > height && Math.random() > 0.975) {
				                arr[i] = 0;
				            }
		            arr[i]++;
		        }
}

setInterval(draw, 30);
window.addEventListener("resize", () => location.reload());
