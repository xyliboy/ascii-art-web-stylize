(function () {
  const rootElement = document.documentElement;

  function updateBackgroundMotion(event) {
    const centerX = window.innerWidth / 2;
    const centerY = window.innerHeight / 2;
    const shiftX = ((event.clientX - centerX) / centerX) * 18;
    const shiftY = ((event.clientY - centerY) / centerY) * 18;

    rootElement.style.setProperty("--star-shift-x", shiftX.toFixed(2) + "px");
    rootElement.style.setProperty("--star-shift-y", shiftY.toFixed(2) + "px");
    rootElement.style.setProperty("--star-counter-x", (shiftX * -0.55).toFixed(2) + "px");
    rootElement.style.setProperty("--star-counter-y", (shiftY * -0.55).toFixed(2) + "px");
  }

  window.addEventListener("pointermove", updateBackgroundMotion);

  const form = document.querySelector(".art-form");
  const textArea = document.getElementById("text");
  const resultElement = document.querySelector(".message-success pre");
  const feedbackElement = document.querySelector("[data-action-feedback]");
  const formFeedbackElement = document.querySelector("[data-form-feedback]");
  const charCountElement = document.querySelector("[data-char-count]");
  const resultFrame = document.querySelector(".result-frame");
  const scrollControlX = document.querySelector('[data-scroll-control="x"]');
  const scrollControlY = document.querySelector('[data-scroll-control="y"]');
  const motionControls = document.querySelectorAll(".scroll-control");
  const copyButton = document.querySelector('[data-action="copy-result"]');
  const saveButton = document.querySelector('[data-action="save-image"]');
  const shareButton = document.querySelector('[data-action="share-result"]');
  const clearButton = document.querySelector('[data-action="clear-form"]');
  const zoomInButton = document.querySelector('[data-action="zoom-in"]');
  const zoomOutButton = document.querySelector('[data-action="zoom-out"]');
  const guidePanel = document.querySelector("[data-guide-panel]");
  const openGuideButton = document.querySelector('[data-action="open-guide"]');
  const closeGuideButton = document.querySelector('[data-action="close-guide"]');
  let resultZoom = 0.95;
  let isDraggingPreview = false;
  let dragStartX = 0;
  let dragStartY = 0;
  let dragStartScrollLeft = 0;
  let dragStartScrollTop = 0;

  if (!form || !textArea) {
    return;
  }

  function setFeedback(message) {
    if (feedbackElement) {
      feedbackElement.textContent = message;
    }
  }

  function setFormFeedback(message) {
    if (formFeedbackElement) {
      formFeedbackElement.textContent = message;
    }
  }

  function updateCharacterCount() {
    if (charCountElement) {
      charCountElement.textContent = String(textArea.value.length);
    }
  }

  function resizeTextArea() {
    textArea.style.height = "auto";
    textArea.style.height = Math.min(textArea.scrollHeight, 520) + "px";
  }

  function openGuide() {
    if (guidePanel) {
      guidePanel.hidden = false;
      guidePanel.scrollIntoView({ behavior: "smooth", block: "start" });
    }
  }

  function closeGuide() {
    if (guidePanel) {
      guidePanel.hidden = true;
    }
    if (openGuideButton) {
      openGuideButton.focus();
    }
  }

  function updateResultZoom() {
    if (!resultElement) {
      return;
    }
    resultElement.style.fontSize = resultZoom.toFixed(2) + "rem";
    setFeedback("Preview zoom: " + Math.round(resultZoom / 0.95 * 100) + "%.");
  }

  function changeResultZoom(delta) {
    resultZoom = Math.min(1.8, Math.max(0.18, resultZoom + delta));
    updateResultZoom();
    updateScrollControls();
  }

  function updateScrollControls() {
    if (!resultFrame) {
      return;
    }

    const maxScrollLeft = Math.max(0, resultFrame.scrollWidth - resultFrame.clientWidth);
    const maxScrollTop = Math.max(0, resultFrame.scrollHeight - resultFrame.clientHeight);

    if (scrollControlX) {
      scrollControlX.max = String(maxScrollLeft);
      scrollControlX.value = String(Math.min(resultFrame.scrollLeft, maxScrollLeft));
      scrollControlX.disabled = maxScrollLeft === 0;
    }

    if (scrollControlY) {
      scrollControlY.max = String(maxScrollTop);
      scrollControlY.value = String(Math.min(resultFrame.scrollTop, maxScrollTop));
      scrollControlY.disabled = maxScrollTop === 0;
      scrollControlY.style.width = resultFrame.clientHeight + "px";
    }
  }

  function syncScrollControls() {
    if (scrollControlX) {
      scrollControlX.value = String(resultFrame ? resultFrame.scrollLeft : 0);
    }
    if (scrollControlY) {
      scrollControlY.value = String(resultFrame ? resultFrame.scrollTop : 0);
    }
  }

  function showControlTrack(control) {
    control.classList.add("is-moving");
    window.clearTimeout(control.motionTimer);
    control.motionTimer = window.setTimeout(function () {
      control.classList.remove("is-moving");
    }, 750);
  }

  function setupMotionControl(control) {
    control.addEventListener("pointerenter", function () {
      showControlTrack(control);
    });
    control.addEventListener("pointermove", function () {
      showControlTrack(control);
    });
    control.addEventListener("input", function () {
      showControlTrack(control);
    });
  }

  function startPreviewDrag(event) {
    if (!resultFrame) {
      return;
    }
    isDraggingPreview = true;
    resultFrame.classList.add("is-dragging");
    dragStartX = event.clientX;
    dragStartY = event.clientY;
    dragStartScrollLeft = resultFrame.scrollLeft;
    dragStartScrollTop = resultFrame.scrollTop;
  }

  function movePreviewDrag(event) {
    if (!isDraggingPreview || !resultFrame) {
      return;
    }
    event.preventDefault();
    resultFrame.scrollLeft = dragStartScrollLeft - (event.clientX - dragStartX);
    resultFrame.scrollTop = dragStartScrollTop - (event.clientY - dragStartY);
    syncScrollControls();
  }

  function stopPreviewDrag() {
    if (!resultFrame) {
      return;
    }
    isDraggingPreview = false;
    resultFrame.classList.remove("is-dragging");
  }

  async function copyResult() {
    if (!resultElement || !navigator.clipboard) {
      setFeedback("Copy is not available in this browser.");
      return;
    }

    await navigator.clipboard.writeText(resultElement.textContent);
    setFeedback("ASCII output copied to the clipboard.");
  }

  function saveResultAsImage() {
    if (!resultElement) {
      setFeedback("There is no ASCII output to export.");
      return;
    }

    const asciiText = resultElement.textContent || "";
    const lines = asciiText.split("\n");
    const fontSize = 16;
    const lineHeight = 20;
    const padding = 28;
    const longestLine = lines.reduce(function (max, line) {
      return Math.max(max, line.length);
    }, 0);

    const canvas = document.createElement("canvas");
    const context = canvas.getContext("2d");
    if (!context) {
      setFeedback("Image export is not available in this browser.");
      return;
    }

    context.font = fontSize + 'px "Courier New", monospace';
    const textWidth = Math.ceil(context.measureText("M").width * Math.max(longestLine, 1));
    canvas.width = textWidth + padding * 2;
    canvas.height = Math.max(lines.length, 1) * lineHeight + padding * 2;

    context.fillStyle = "#ffffff";
    context.fillRect(0, 0, canvas.width, canvas.height);
    context.fillStyle = "#172033";
    context.font = fontSize + 'px "Courier New", monospace';
    context.textBaseline = "top";

    lines.forEach(function (line, index) {
      context.fillText(line, padding, padding + index * lineHeight);
    });

    const link = document.createElement("a");
    link.href = canvas.toDataURL("image/png");
    link.download = "ascii-art.png";
    link.click();
    setFeedback("ASCII output saved as ascii-art.png.");
  }

  async function shareResult() {
    if (!resultElement) {
      setFeedback("There is no ASCII output to share.");
      return;
    }

    const shareText = resultElement.textContent || "";
    if (navigator.share) {
      await navigator.share({
        title: "ASCII Art Web Output",
        text: shareText
      });
      setFeedback("Share dialog opened.");
      return;
    }

    if (navigator.clipboard) {
      await navigator.clipboard.writeText(shareText);
      setFeedback("Native share is not available, so the output was copied instead.");
      return;
    }

    setFeedback("Share is not available in this browser.");
  }

  function clearForm() {
    textArea.value = "";
    updateCharacterCount();
    resizeTextArea();
    setFormFeedback("Text cleared. You can start again.");
    textArea.focus();
  }

  updateCharacterCount();
  resizeTextArea();
  updateScrollControls();
  motionControls.forEach(setupMotionControl);
  textArea.addEventListener("input", function () {
    updateCharacterCount();
    resizeTextArea();
    if (textArea.value.length === 0) {
      setFormFeedback("The text box is empty.");
      return;
    }
    setFormFeedback("Ready to make ASCII art.");
  });

  textArea.addEventListener("keydown", function (event) {
    if (event.key !== "Enter" || event.shiftKey) {
      return;
    }

    event.preventDefault();
    if (typeof form.requestSubmit === "function") {
      form.requestSubmit();
      return;
    }

    form.submit();
  });

  if (copyButton) {
    copyButton.addEventListener("click", function () {
      copyResult().catch(function () {
        setFeedback("Copy failed. Please try again.");
      });
    });
  }

  if (saveButton) {
    saveButton.addEventListener("click", function () {
      saveResultAsImage();
    });
  }

  if (shareButton) {
    shareButton.addEventListener("click", function () {
      shareResult().catch(function () {
        setFeedback("Share failed. Please try again.");
      });
    });
  }

  if (clearButton) {
    clearButton.addEventListener("click", function () {
      clearForm();
    });
  }

  if (zoomInButton) {
    zoomInButton.addEventListener("click", function () {
      changeResultZoom(0.1);
    });
  }

  if (zoomOutButton) {
    zoomOutButton.addEventListener("click", function () {
      changeResultZoom(-0.1);
    });
  }

  if (resultFrame) {
    resultFrame.addEventListener("mousedown", startPreviewDrag);
    resultFrame.addEventListener("mousemove", movePreviewDrag);
    resultFrame.addEventListener("mouseup", stopPreviewDrag);
    resultFrame.addEventListener("mouseleave", stopPreviewDrag);
    resultFrame.addEventListener("scroll", syncScrollControls);
    window.addEventListener("mouseup", stopPreviewDrag);
    window.addEventListener("resize", updateScrollControls);
  }

  if (scrollControlX) {
    scrollControlX.addEventListener("input", function () {
      if (resultFrame) {
        resultFrame.scrollLeft = Number(scrollControlX.value);
      }
    });
  }

  if (scrollControlY) {
    scrollControlY.addEventListener("input", function () {
      if (resultFrame) {
        resultFrame.scrollTop = Number(scrollControlY.value);
      }
    });
  }

  if (openGuideButton) {
    openGuideButton.addEventListener("click", function () {
      openGuide();
    });
  }

  if (closeGuideButton) {
    closeGuideButton.addEventListener("click", function () {
      closeGuide();
    });
  }
}());
