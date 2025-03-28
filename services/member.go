package services

import (
	"errors"
	"huiyuanguanli/dto"
	"huiyuanguanli/models"
	"time"

	"gorm.io/gorm"
)

// MemberService 会员服务
type MemberService struct {
	db *gorm.DB
}

// NewMemberService 创建会员服务
func NewMemberService(db *gorm.DB) *MemberService {
	return &MemberService{db: db}
}

// CreateMember 创建会员
func (s *MemberService) CreateMember(req *dto.CreateMemberRequest) (*models.Member, error) {
	// 检查手机号是否已存在
	var count int64
	if err := s.db.Model(&models.Member{}).Where("phone = ?", req.Phone).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("手机号已存在")
	}

	// 检查会员卡号是否已存在
	if err := s.db.Model(&models.Member{}).Where("card_id = ?", req.CardID).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("会员卡号已存在")
	}

	// 创建会员
	member := &models.Member{
		Name:    req.Name,
		Gender:  req.Gender,
		Phone:   req.Phone,
		CardID:  req.CardID,
		StoreID: req.StoreID,
		Balance: req.Balance,
		Remark:  req.Remark,
	}

	// 开启事务
	tx := s.db.Begin()
	if err := tx.Create(member).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// 添加标签
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		if err := tx.Find(&tags, req.TagIDs).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
		if err := tx.Model(member).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	// 如果有初始余额，创建充值记录
	if req.Balance > 0 {
		transaction := &models.Transaction{
			MemberID:    member.ID,
			StoreID:     req.StoreID,
			Type:        "RECHARGE",
			Amount:      req.Balance,
			Balance:     req.Balance,
			GiftBalance: 0,
			Description: "初始充值",
		}
		if err := tx.Create(transaction).Error; err != nil {
			tx.Rollback()
			return nil, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return member, nil
}

// UpdateMember 更新会员
func (s *MemberService) UpdateMember(id uint, req *dto.UpdateMemberRequest) error {
	// 检查会员是否存在
	var member models.Member
	if err := s.db.First(&member, id).Error; err != nil {
		return err
	}

	// 检查手机号是否已被其他会员使用
	var count int64
	if err := s.db.Model(&models.Member{}).Where("phone = ? AND id != ?", req.Phone, id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("手机号已被其他会员使用")
	}

	// 检查会员卡号是否已被其他会员使用
	if err := s.db.Model(&models.Member{}).Where("card_id = ? AND id != ?", req.CardID, id).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return errors.New("会员卡号已被其他会员使用")
	}

	// 开启事务
	tx := s.db.Begin()

	// 更新会员信息
	updates := map[string]interface{}{
		"name":     req.Name,
		"gender":   req.Gender,
		"phone":    req.Phone,
		"card_id":  req.CardID,
		"store_id": req.StoreID,
		"remark":   req.Remark,
	}
	if err := tx.Model(&member).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 更新标签
	if len(req.TagIDs) > 0 {
		var tags []models.Tag
		if err := tx.Find(&tags, req.TagIDs).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Model(&member).Association("Tags").Replace(tags); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

// GetMemberList 获取会员列表
func (s *MemberService) GetMemberList(req *dto.MemberListRequest) (*dto.MemberListResponse, error) {
	var response dto.MemberListResponse
	query := s.db.Model(&models.Member{})

	// 添加查询条件
	if req.Name != "" {
		query = query.Where("name LIKE ?", "%"+req.Name+"%")
	}
	if req.Phone != "" {
		query = query.Where("phone LIKE ?", "%"+req.Phone+"%")
	}
	if req.CardID != "" {
		query = query.Where("card_id LIKE ?", "%"+req.CardID+"%")
	}
	if req.StoreID > 0 {
		query = query.Where("store_id = ?", req.StoreID)
	}
	if req.TagID > 0 {
		query = query.Joins("JOIN member_tags ON members.id = member_tags.member_id").
			Where("member_tags.tag_id = ?", req.TagID)
	}
	if req.Status >= 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 获取总数
	if err := query.Count(&response.Total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	var members []models.Member
	if err := query.Preload("Store").
		Preload("Tags").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&members).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	response.Items = make([]dto.MemberDetail, len(members))
	for i, member := range members {
		response.Items[i] = dto.MemberDetail{
			ID:              member.ID,
			Name:            member.Name,
			Gender:          member.Gender,
			Phone:           member.Phone,
			CardID:          member.CardID,
			StoreID:         member.StoreID,
			StoreName:       member.Store.Name,
			Balance:         member.Balance,
			GiftBalance:     member.GiftBalance,
			TotalRecharge:   member.TotalRecharge,
			TotalConsume:    member.TotalConsume,
			FirstActiveTime: member.FirstActiveTime.Format("2006-01-02 15:04:05"),
			LastActiveTime:  member.LastActiveTime.Format("2006-01-02 15:04:05"),
			Status:          member.Status,
			Remark:          member.Remark,
			CreatedAt:       member.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:       member.UpdatedAt.Format("2006-01-02 15:04:05"),
		}

		// 转换标签
		response.Items[i].Tags = make([]dto.Tag, len(member.Tags))
		for j, tag := range member.Tags {
			response.Items[i].Tags[j] = dto.Tag{
				ID:   tag.ID,
				Name: tag.Name,
			}
		}
	}

	return &response, nil
}

// GetMemberDetail 获取会员详情
func (s *MemberService) GetMemberDetail(id uint) (*dto.MemberDetail, error) {
	var member models.Member
	if err := s.db.Preload("Store").
		Preload("Tags").
		First(&member, id).Error; err != nil {
		return nil, err
	}

	detail := &dto.MemberDetail{
		ID:              member.ID,
		Name:            member.Name,
		Gender:          member.Gender,
		Phone:           member.Phone,
		CardID:          member.CardID,
		StoreID:         member.StoreID,
		StoreName:       member.Store.Name,
		Balance:         member.Balance,
		GiftBalance:     member.GiftBalance,
		TotalRecharge:   member.TotalRecharge,
		TotalConsume:    member.TotalConsume,
		FirstActiveTime: member.FirstActiveTime.Format("2006-01-02 15:04:05"),
		LastActiveTime:  member.LastActiveTime.Format("2006-01-02 15:04:05"),
		Status:          member.Status,
		Remark:          member.Remark,
		CreatedAt:       member.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       member.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	// 转换标签
	detail.Tags = make([]dto.Tag, len(member.Tags))
	for i, tag := range member.Tags {
		detail.Tags[i] = dto.Tag{
			ID:   tag.ID,
			Name: tag.Name,
		}
	}

	return detail, nil
}

// DeleteMember 删除会员
func (s *MemberService) DeleteMember(id uint) error {
	return s.db.Model(&models.Member{}).Where("id = ?", id).Update("status", 0).Error
}

// Recharge 会员充值
func (s *MemberService) Recharge(id uint, req *dto.RechargeRequest, operatorID uint) error {
	// 检查会员是否存在
	var member models.Member
	if err := s.db.First(&member, id).Error; err != nil {
		return err
	}

	// 开启事务
	tx := s.db.Begin()

	// 更新会员余额
	updates := map[string]interface{}{
		"balance":          member.Balance + req.Amount,
		"gift_balance":     member.GiftBalance + req.GiftAmount,
		"total_recharge":   member.TotalRecharge + req.Amount,
		"last_active_time": time.Now(),
	}
	if err := tx.Model(&member).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建充值记录
	transaction := &models.Transaction{
		MemberID:    member.ID,
		StoreID:     member.StoreID,
		Type:        "RECHARGE",
		Amount:      req.Amount,
		Balance:     member.Balance + req.Amount,
		GiftBalance: member.GiftBalance + req.GiftAmount,
		Description: req.Description,
		OperatorID:  operatorID,
	}
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Consume 会员消费
func (s *MemberService) Consume(id uint, req *dto.ConsumeRequest, operatorID uint) error {
	// 检查会员是否存在
	var member models.Member
	if err := s.db.First(&member, id).Error; err != nil {
		return err
	}

	// 检查余额是否足够
	if member.Balance < req.Amount {
		return errors.New("余额不足")
	}

	// 开启事务
	tx := s.db.Begin()

	// 更新会员余额
	updates := map[string]interface{}{
		"balance":          member.Balance - req.Amount,
		"total_consume":    member.TotalConsume + req.Amount,
		"last_active_time": time.Now(),
	}
	if err := tx.Model(&member).Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 创建消费记录
	transaction := &models.Transaction{
		MemberID:    member.ID,
		StoreID:     member.StoreID,
		Type:        "CONSUME",
		Amount:      req.Amount,
		Balance:     member.Balance - req.Amount,
		GiftBalance: member.GiftBalance,
		ItemType:    req.ItemType,
		ItemID:      req.ItemID,
		Description: req.Description,
		OperatorID:  operatorID,
	}
	if err := tx.Create(transaction).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetTransactionList 获取交易记录列表
func (s *MemberService) GetTransactionList(req *dto.TransactionListRequest) (*dto.TransactionListResponse, error) {
	var response dto.TransactionListResponse
	query := s.db.Model(&models.Transaction{}).Where("member_id = ?", req.MemberID)

	// 添加查询条件
	if req.Type != "" {
		query = query.Where("type = ?", req.Type)
	}
	if req.StartTime != "" {
		startTime, err := time.Parse("2006-01-02", req.StartTime)
		if err != nil {
			return nil, err
		}
		query = query.Where("created_at >= ?", startTime)
	}
	if req.EndTime != "" {
		endTime, err := time.Parse("2006-01-02", req.EndTime)
		if err != nil {
			return nil, err
		}
		query = query.Where("created_at <= ?", endTime)
	}

	// 获取总数
	if err := query.Count(&response.Total).Error; err != nil {
		return nil, err
	}

	// 获取分页数据
	var transactions []models.Transaction
	if err := query.Preload("Store").
		Preload("Operator").
		Order("created_at DESC").
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		Find(&transactions).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	response.Items = make([]dto.TransactionDetail, len(transactions))
	for i, transaction := range transactions {
		response.Items[i] = dto.TransactionDetail{
			ID:           transaction.ID,
			Type:         transaction.Type,
			Amount:       transaction.Amount,
			Balance:      transaction.Balance,
			GiftBalance:  transaction.GiftBalance,
			ItemType:     transaction.ItemType,
			ItemID:       transaction.ItemID,
			Description:  transaction.Description,
			StoreID:      transaction.StoreID,
			StoreName:    transaction.Store.Name,
			OperatorID:   transaction.OperatorID,
			OperatorName: transaction.Operator.Name,
			CreatedAt:    transaction.CreatedAt.Format("2006-01-02 15:04:05"),
		}
	}

	return &response, nil
}
