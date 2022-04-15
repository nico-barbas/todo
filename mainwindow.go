package main

import (
	"image/color"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	maxSessionPerLine   = 10
	sessionCheckPadding = 4
	sessionCheckHeight  = 600 / (maxSessionPerLine + sessionCheckPadding*2)
	sessionCheckSize    = sessionCheckHeight - (sessionCheckPadding * 2)
)

const (
	settingsBtnID rectID = iota
	archiveBtnID
	timerBtnID
	taskSettingsBtnID
	archiveTaskBtnID
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

	timerBtnPressed bool

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

func (m *mainWindow) update(mPos point, mLeft bool, selected bool) (startTask bool) {
	m.timerBtnPressed = false
	if !isInputHandled(mPos) {
		if selected {
			m.infoElements.update(mPos, mLeft)
		}
		m.settingElements.update(mPos, mLeft)
	}
	return m.timerBtnPressed
}

func (m *mainWindow) draw(dst *ebiten.Image, task *task) {
	ebitenutil.DrawLine(
		dst,
		m.rect.full.x,
		m.rect.full.y,
		m.rect.full.x,
		m.rect.full.y+m.rect.full.height,
		color.White,
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
		drawRect(dst, m.progressRect.remaining, White)

		drawTextBtn(dst, m.workTimerRect.remaining, string(task.timer.toString()), largeTextSize)

		drawTextBtn(dst, m.restTimerRect.remaining, string(task.timer.toString()), largeTextSize)

		drawTextBtn(dst, m.timerBtnRect.remaining, "Start Timer", textSize)

		drawTextBtn(dst, m.archiveTaskBtnRect.remaining, "Archive Task", textSize)
		drawIcontBtn(dst, m.taskSettingsBtnRect.remaining, m.archiveIcon)
	}
}

func (m *mainWindow) onClick(userID rectID) {
	switch userID {
	case timerBtnID:
		m.timerBtnPressed = true
	case archiveTaskBtnID:
		FireSignal(todoTaskRemoved, SignalNoArgs)
	}
}
