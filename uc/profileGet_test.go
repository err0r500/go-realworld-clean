package uc_test

//
// func TestProfileGet_happyCase(t *testing.T) {
// 	t.Run("authenticated", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()
//
// 		rick := testData.User("rick")
// 		jane := testData.User("jane")
// 		rick.FollowIDs = []string{jane.Name}
//
// 		i := mock.NewMockedInteractor(mockCtrl)
// 		i.UserRW.EXPECT().GetByName(rick.Name).Return(&rick, nil).Times(1)
// 		i.UserRW.EXPECT().GetByName(jane.Name).Return(&jane, nil).Times(1)
//
// 		retProfile, follows, err := i.GetUCHandler().ProfileGet(rick.Name, jane.Name)
//
// 		assert.NoError(t, err)
// 		assert.True(t, follows)
// 		assert.Equal(t, jane, *retProfile)
// 	})
//
// 	t.Run("not authenticated", func(t *testing.T) {
// 		mockCtrl := gomock.NewController(t)
// 		defer mockCtrl.Finish()
//
// 		jane := testData.User("jane")
//
// 		i := mock.NewMockedInteractor(mockCtrl)
// 		i.UserRW.EXPECT().GetByName(jane.Name).Return(&jane, nil).Times(1)
//
// 		retProfile, follows, err := i.GetUCHandler().ProfileGet("", jane.Name)
//
// 		assert.NoError(t, err)
// 		assert.False(t, follows)
// 		assert.Equal(t, jane, *retProfile)
// 	})
// }
