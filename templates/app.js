(function () {
  const textarea = document.getElementById("text");

  if (textarea) {
    const hintText = "Write anything... then pick a banner ";
    let i = 0;

    function typeWriter() {
      textarea.placeholder = hintText.slice(0, i + 1);
      i++;
      if (i < hintText.length) {
        setTimeout(typeWriter, 150);
      } else {
        setTimeout(function () {
          i = 0;
          typeWriter();
        }, 2500);
      }
    }

    typeWriter();

    textarea.addEventListener("keydown", function (e) {
      if (e.key === "Enter" && !e.shiftKey) {
        e.preventDefault();
        textarea.closest("form").submit();
      }
    });
  }

  const bannerInputs = Array.from(document.querySelectorAll('input[name="banner"]'));
  document.addEventListener("keydown", function (e) {
    if (e.target === textarea) return;
    const current = bannerInputs.findIndex(function (input) { return input.checked; });
    if (e.key === "ArrowRight" && current < bannerInputs.length - 1) {
      bannerInputs[current + 1].checked = true;
      bannerInputs[current + 1].focus();
    }
    if (e.key === "ArrowLeft" && current > 0) {
      bannerInputs[current - 1].checked = true;
      bannerInputs[current - 1].focus();
    }
  });

  const output = document.getElementById("ascii-output");
  const copyBtn = document.getElementById("btn-copy");
  const saveBtn = document.getElementById("btn-save");
  const shareBtn = document.getElementById("btn-share");
  const feedback = document.getElementById("action-feedback");

  if (!output) {
    return;
  }

  function setFeedback(message) {
    feedback.textContent = message;
    setTimeout(function () {
      feedback.textContent = "";
    }, 3000);
  }

  if (copyBtn) {
    copyBtn.addEventListener("click", function () {
      navigator.clipboard.writeText(output.textContent).then(function () {
        setFeedback("Copied to clipboard.");
      }).catch(function () {
        setFeedback("Copy failed. Please try again.");
      });
    });
  }

  if (saveBtn) {
    saveBtn.addEventListener("click", function () {
      const text = output.textContent || "";
      const lines = text.split("\n");
      const fontSize = 16;
      const lineHeight = 20;
      const padding = 24;

      const canvas = document.createElement("canvas");
      const ctx = canvas.getContext("2d");
      ctx.font = fontSize + "px 'Courier New', monospace";

      const longestLine = lines.reduce(function (max, line) {
        return Math.max(max, line.length);
      }, 0);

      canvas.width = Math.ceil(ctx.measureText("M").width * longestLine) + padding * 2;
      canvas.height = lines.length * lineHeight + padding * 2;

      ctx.fillStyle = "#ffffff";
      ctx.fillRect(0, 0, canvas.width, canvas.height);
      ctx.fillStyle = "#172033";
      ctx.font = fontSize + "px 'Courier New', monospace";
      ctx.textBaseline = "top";

      lines.forEach(function (line, i) {
        ctx.fillText(line, padding, padding + i * lineHeight);
      });

      const link = document.createElement("a");
      link.href = canvas.toDataURL("image/png");
      link.download = "ascii-art.png";
      link.click();
      setFeedback("Image saved as ascii-art.png.");
    });
  }

  if (shareBtn) {
    shareBtn.addEventListener("click", function () {
      const text = output.textContent || "";
      if (navigator.share) {
        navigator.share({ title: "ASCII Art", text: text }).then(function () {
          setFeedback("Shared successfully.");
        }).catch(function () {
          setFeedback("Share cancelled.");
        });
      } else if (navigator.clipboard) {
        navigator.clipboard.writeText(text).then(function () {
          setFeedback("Share not available. Copied to clipboard instead.");
        });
      } else {
        setFeedback("Share is not available in this browser.");
      }
    });
  }
}());
