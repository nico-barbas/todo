package main

import (
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	settingsBtnID rectID = iota
	archiveBtnID
	timerBtnID
	taskSettingsBtnID
	archiveTaskBtnID
)

const (
	timerStartStr = "Start"
	timerStopStr  = "Stop"
)

type mainWindow struct {
	rect                rectLayout
	settingsBtnRect     rectLayout
	archiveBtnRect      rectLayout
	titleRect           rectLayout
	progressRect        rectLayout
	workTimerRect       rectLayout
	restTimerRect       rectLayout
	timerBtnRect        rectLayout
	taskSettingsBtnRect rectLayout
	archiveTaskBtnRect  rectLayout

	settingElements rectArray
	infoElements    rectArray

	timerStr string

	font          *Font
	rectOutline   *ebiten.Image
	outlineConstr constraint
	archiveIcon   *ebiten.Image
	settingsIcon  *ebiten.Image
}

func (m *mainWindow) init(font *Font, outline *ebiten.Image) {
	const mainWindowPadding = 10
	const mainWindowNoPadding = 0
	m.settingElements.init(m, 5)
	m.infoElements.init(m, 5)

	m.rect = newRectLayout(rectangle{200, 0, 600, windowHeight})
	m.rect.cut(rectCutUp, mainWindowPadding, 0)

	settingsRect := m.rect.cut(rectCutUp, 30, mainWindowPadding)
	settingsRect.cut(rectCutRight, mainWindowPadding, 0)
	m.settingsBtnRect = settingsRect.cut(rectCutRight, 30, mainWindowPadding)
	m.archiveBtnRect = settingsRect.cut(rectCutRight, 30, mainWindowPadding)
	m.settingElements.add(m.settingsBtnRect.remaining, settingsBtnID)
	m.settingElements.add(m.archiveBtnRect.remaining, archiveBtnID)

	m.titleRect = m.rect.cut(rectCutUp, 80, mainWindowPadding+10)

	m.progressRect = m.rect.cut(rectCutUp, 50, mainWindowPadding)
	m.progressRect.cut(rectCutRight, 200, 0)
	m.progressRect.cut(rectCutLeft, 200, 0)

	timerRect := m.rect.cut(rectCutUp, 100, mainWindowPadding)
	timerRect.cut(rectCutRight, 100, 0)
	timerRect.cut(rectCutLeft, 100, 0)
	timerWidth := (timerRect.remaining.width - 5) / 2
	m.workTimerRect = timerRect.cut(rectCutLeft, timerWidth, mainWindowNoPadding)
	m.restTimerRect = timerRect.cut(rectCutRight, timerWidth, mainWindowNoPadding)

	m.timerBtnRect = m.rect.cut(rectCutUp, 50, mainWindowPadding)
	m.timerBtnRect.cut(rectCutRight, 200, 0)
	m.timerBtnRect.cut(rectCutLeft, 200, 0)
	m.infoElements.add(m.timerBtnRect.remaining, timerBtnID)

	m.rect.cut(rectCutDown, mainWindowPadding, 0)
	taskSettingsRect := m.rect.cut(rectCutDown, 30, mainWindowPadding)
	taskSettingsRect.cut(rectCutRight, mainWindowPadding, 0)
	m.archiveTaskBtnRect = taskSettingsRect.cut(
		rectCutRight,
		font.MeasureText("Archive Task", textSize)[0]+mainWindowPadding*2,
		mainWindowPadding,
	)
	m.taskSettingsBtnRect = taskSettingsRect.cut(rectCutRight, 30, mainWindowPadding)
	m.infoElements.add(m.archiveTaskBtnRect.remaining, archiveTaskBtnID)
	m.infoElements.add(m.taskSettingsBtnRect.remaining, taskSettingsBtnID)

	m.font = font
	m.rectOutline = outline
	m.outlineConstr = constraint{2, 2, 2, 2}

	m.archiveIcon, _, _ = ebitenutil.NewImageFromFile("assets/icon-archive.png")
	m.settingsIcon, _, _ = ebitenutil.NewImageFromFile("assets/icon-settings.png")
}

func (m *mainWindow) update(mPos point, mLeft bool, task *task) {
	if !isInputHandled(mPos) {
		if task != nil {
			m.infoElements.update(mPos, mLeft)
			switch task.timer.running {
			case true:
				m.timerStr = timerStopStr
			case false:
				m.timerStr = timerStartStr
			}
		}
		m.settingElements.update(mPos, mLeft)
	}
}

func (m *mainWindow) draw(dst *ebiten.Image, task *task) {
	drawRect(dst, m.rect.full, darkBackground2)
	ebitenutil.DrawLine(
		dst,
		m.rect.full.x,
		m.rect.full.y,
		m.rect.full.x,
		m.rect.full.y+m.rect.full.height,
		darkSeparator,
	)
	m.settingElements.highlight(dst)

	drawIcontBtn(dst, m.settingsBtnRect.remaining, m.settingsIcon)
	drawIcontBtn(dst, m.archiveBtnRect.remaining, m.archiveIcon)

	// Task info widgets
	if task != nil {
		m.infoElements.highlight(dst)

		drawTextCenter(dst, textOptions{
			font: m.font, text: task.name, bounds: m.titleRect.remaining,
			size: largeTextSize, clr: White,
		})
		// drawRect(dst, m.progressRect.remaining, White)
		drawImageSlice(dst, m.progressRect.remaining, rectOutline, rectConstraint, White)
		// Draw the progress bars here
		{
			insideRect := m.progressRect.remaining
			insideRect.x += 3
			insideRect.y += 3
			insideRect.width -= 4
			insideRect.height -= 6
			barWidth := (insideRect.width - float64(2*task.sessionRequired)) / float64(task.sessionRequired)
			xptr := insideRect.x
			for i := 0; i < task.sessionCompleted; i += 1 {
				drawRect(dst, rectangle{xptr, insideRect.y, barWidth, insideRect.height}, White)
				xptr += barWidth + 2
			}
		}

		if task.isWorkInProgress() {
			progress := task.progress()
			drawRect(dst, rectangle{
				m.workTimerRect.remaining.x, m.workTimerRect.remaining.y,
				m.workTimerRect.remaining.width * progress, m.workTimerRect.remaining.height,
			}, WhiteA125)
		}
		drawTextBtn(dst, m.workTimerRect.remaining, task.getWorkTime(), largeTextSize)
		drawImage(
			dst, timerWorkIcon,
			point{
				m.workTimerRect.x() + m.workTimerRect.width()/2 - 8,
				m.workTimerRect.y() + 8,
			},
			WhiteA125,
		)

		if task.isRestInProgress() {
			progress := task.progress()
			drawRect(dst, rectangle{
				m.restTimerRect.remaining.x, m.restTimerRect.remaining.y,
				m.restTimerRect.remaining.width * progress, m.restTimerRect.remaining.height,
			}, WhiteA125)
		}
		drawTextBtn(dst, m.restTimerRect.remaining, task.getRestTime(), largeTextSize)
		drawImage(
			dst, timerRestIcon,
			point{
				m.restTimerRect.x() + m.restTimerRect.width()/2 - 8,
				m.restTimerRect.y() + 8,
			},
			WhiteA125,
		)

		// Could probably cache this string
		// maybe no allocations are even happening.. who knows
		drawTextBtn(dst, m.timerBtnRect.remaining, m.timerStr+" Timer", textSize)

		drawTextBtn(dst, m.archiveTaskBtnRect.remaining, "Archive Task", textSize)
		drawIcontBtn(dst, m.taskSettingsBtnRect.remaining, m.archiveIcon)
	}
}

func (m *mainWindow) onClick(userID rectID) {
	switch userID {
	case archiveBtnID:
		FireSignal(todoArchiveBtnPressed, SignalNoArgs)
	case timerBtnID:
		switch m.timerStr {
		case timerStartStr:
			FireSignal(todoTaskStarted, SignalNoArgs)
		case timerStopStr:
			FireSignal(todoTaskStopped, SignalNoArgs)
		}
	case archiveTaskBtnID:
		FireSignal(todoTaskRemoved, SignalNoArgs)
	}
}
