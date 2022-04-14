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
	settingsBtnRect rectID = iota
	archiveBtnRect
	titleRect
	progressRect
	workTimerRect
	restTimerRect
	timerBtnRect
	taskSettingsBtnRect
	archiveTaskBtnRect
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

	font          *Font
	rectOutline   *ebiten.Image
	outlineConstr constraint
	archiveIcon   *ebiten.Image
	settingsIcon  *ebiten.Image

	shouldHighlight bool
	highlightRect   rectangle
}

func (m *mainWindow) init(font *Font, outline *ebiten.Image) {
	const mainWindowPadding = 10
	const mainWindowNoPadding = 0
	m.settingElements.init(m, 5)
	m.infoElements.init(m, 20)

	m.rect = newRectLayout(rectangle{200, 0, 600, windowHeight})
	m.rect.cut(rectCutUp, mainWindowPadding, 0)

	settingsRect := m.rect.cut(rectCutUp, 30, mainWindowPadding)
	settingsRect.cut(rectCutRight, mainWindowPadding, 0)
	m.settingsBtnRect = settingsRect.cut(rectCutRight, 30, mainWindowPadding)
	m.archiveBtnRect = settingsRect.cut(rectCutRight, 30, mainWindowPadding)

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

	m.rect.cut(rectCutDown, mainWindowPadding, 0)
	taskSettingsRect := m.rect.cut(rectCutDown, 30, mainWindowPadding)
	taskSettingsRect.cut(rectCutRight, mainWindowPadding, 0)
	m.archiveTaskBtnRect = taskSettingsRect.cut(
		rectCutRight,
		font.MeasureText("Archive Task", textSize)[0]+mainWindowPadding*2,
		mainWindowPadding,
	)
	m.taskSettingsBtnRect = taskSettingsRect.cut(rectCutRight, 30, mainWindowPadding)

	m.font = font
	m.rectOutline = outline
	m.outlineConstr = constraint{2, 2, 2, 2}

	m.archiveIcon, _, _ = ebitenutil.NewImageFromFile("assets/icon-archive.png")
	m.settingsIcon, _, _ = ebitenutil.NewImageFromFile("assets/icon-settings.png")
}

func (m *mainWindow) update(mPos point, mLeft bool, selected bool) (startTask bool) {

	m.shouldHighlight = false
	if !isInputHandled(mPos) {
		if selected {
			switch {
			case m.timerBtnRect.remaining.boundCheck(mPos):
				m.shouldHighlight = true
				m.highlightRect = m.timerBtnRect.remaining
				if mLeft {
					startTask = true
				}
			case m.archiveTaskBtnRect.remaining.boundCheck(mPos):
				m.shouldHighlight = true
				m.highlightRect = m.archiveTaskBtnRect.remaining

			case m.taskSettingsBtnRect.remaining.boundCheck(mPos):
				m.shouldHighlight = true
				m.highlightRect = m.taskSettingsBtnRect.remaining
			}

		}
		switch {
		case m.archiveBtnRect.remaining.boundCheck(mPos):
			m.shouldHighlight = true
			m.highlightRect = m.archiveBtnRect.remaining

		case m.settingsBtnRect.remaining.boundCheck(mPos):
			m.shouldHighlight = true
			m.highlightRect = m.settingsBtnRect.remaining
		}
	}
	return
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
	if m.shouldHighlight {
		drawRect(dst, m.highlightRect, WhiteA125)
	}

	drawImageSlice(dst, m.settingsBtnRect.remaining, m.rectOutline, m.outlineConstr, White)
	drawImageCentered(dst, m.settingsIcon, m.settingsBtnRect.remaining, 1, White)
	drawImageSlice(dst, m.archiveBtnRect.remaining, m.rectOutline, m.outlineConstr, White)
	drawImageCentered(dst, m.archiveIcon, m.archiveBtnRect.remaining, 1, White)

	if task != nil {

		drawTextCenter(dst, textOptions{
			font: m.font, text: task.name, bounds: m.titleRect.remaining,
			size: largeTextSize, clr: White,
		})
		drawRect(dst, m.progressRect.remaining, White)

		drawImageSlice(dst, m.workTimerRect.remaining, m.rectOutline, m.outlineConstr, White)
		drawTextCenter(dst, textOptions{
			font: m.font, text: string(task.timer.toString()), bounds: m.workTimerRect.remaining,
			size: largeTextSize, clr: White,
		})

		drawImageSlice(dst, m.restTimerRect.remaining, m.rectOutline, m.outlineConstr, White)
		drawTextCenter(dst, textOptions{
			font: m.font, text: string(task.timer.toString()), bounds: m.restTimerRect.remaining,
			size: largeTextSize, clr: White,
		})

		drawImageSlice(dst, m.timerBtnRect.remaining, m.rectOutline, m.outlineConstr, White)
		drawTextCenter(dst, textOptions{
			font: m.font, text: "Start Timer", bounds: m.timerBtnRect.remaining,
			size: textSize, clr: White,
		})

		drawImageSlice(dst, m.archiveTaskBtnRect.remaining, m.rectOutline, m.outlineConstr, White)
		drawTextCenter(dst, textOptions{
			font: m.font, text: "Archive Task", bounds: m.archiveTaskBtnRect.remaining,
			size: textSize, clr: White,
		})
		drawImageSlice(dst, m.taskSettingsBtnRect.remaining, m.rectOutline, m.outlineConstr, White)
	}
}

func (m *mainWindow) onClick(userID rectID) {

}
